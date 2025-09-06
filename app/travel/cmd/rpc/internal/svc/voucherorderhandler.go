package svc

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"heart-trip/app/travel/cmd/rpc/internal/config"
	"heart-trip/app/travel/cmd/rpc/pb"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

type VoucherOrderHandler struct {
	cfg         config.Config
	redisClient *redis.Client
	svcCtx      *ServiceContext
}

func NewVoucherOrderHandler(cfg config.Config, svcCtx *ServiceContext, redisClient *redis.Client) *VoucherOrderHandler {
	return &VoucherOrderHandler{
		cfg:         cfg,
		redisClient: redisClient,
		svcCtx:      svcCtx,
	}
}

func (h *VoucherOrderHandler) Start() {
	for {
		// 使用 XREADGROUP 读取消息
		streams, err := h.redisClient.XReadGroup(context.Background(), &redis.XReadGroupArgs{
			Group:    "g1",
			Consumer: "c1",
			Streams:  []string{"stream.orders", ">"},
			Count:    1,
			Block:    2000 * time.Millisecond,
		}).Result()

		if err != nil {
			if errors.Is(err, redis.Nil) {
				// 没有消息，继续循环
				continue
			}
			log.Printf("读取 Redis Stream 失败: %v", err)
			time.Sleep(1 * time.Second)
			continue
		}

		for _, stream := range streams {
			for _, message := range stream.Messages {
				var order pb.VoucherOrder
				err := json.Unmarshal([]byte(message.Values["orderId"].(string)), &order)
				if err != nil {
					log.Printf("解析订单失败: %v", err)
					// 这里可以决定是否要进行 XACK
					continue
				}

				// 处理订单
				h.createVoucherOrder(&order)

				// 确认消息
				ackErr := h.redisClient.XAck(context.Background(), "stream.orders", "g1", message.ID).Err()
				if ackErr != nil {
					log.Printf("确认 Redis 消息失败: %v", ackErr)
				}
			}
		}
	}
}

func (h *VoucherOrderHandler) createVoucherOrder(order *pb.VoucherOrder) {
	userId := order.UserId
	voucherId := order.VoucherId

	// 创建锁对象
	lock := h.svcCtx.RedSync.NewMutex(fmt.Sprintf("lock:order:%d", userId))

	// 尝试获取锁
	if err := lock.Lock(); err != nil {
		log.Printf("获取锁失败: %v", err)
		return
	}

	defer func() {
		if ok, err := lock.Unlock(); !ok || err != nil {
			log.Printf("释放锁失败: %v", err)
		}
	}()

	// 查询订单
	var count int64
	h.svcCtx.DB.Model(&pb.VoucherOrder{}).
		Where("user_id = ? AND voucher_id = ?", userId, voucherId).
		Count(&count)
	if count > 0 {
		log.Printf("用户 %d 已经购买过该 Voucher %d", userId, voucherId)
		return
	}

	// 扣减库存
	tx := h.svcCtx.DB.Begin()
	var stock int64
	tx.Model(&pb.SeckillVoucher{}).Where("id = ?", voucherId).Select("stock").Scan(&stock)
	if stock <= 0 {
		tx.Rollback()
		log.Printf("Voucher %d 库存不足", voucherId)
		return
	}
	tx.Model(&pb.SeckillVoucher{}).Where("id = ?", voucherId).UpdateColumn("stock", gorm.Expr("stock - ?", 1))

	// 创建订单
	tx.Create(&order)

	tx.Commit()
}
