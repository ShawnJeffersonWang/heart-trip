package logic

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/zeromicro/go-zero/core/logx"
	"golodge/app/travel/cmd/rpc/internal/svc"
	"golodge/app/travel/cmd/rpc/pb"
	"golodge/common/ctxdata"
	"golodge/common/tool"
	"os"
	"strconv"
)

type SeckillVoucherLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSeckillVoucherLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SeckillVoucherLogic {
	return &SeckillVoucherLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SeckillVoucherLogic) SeckillVoucher(req *pb.SeckillVoucherRequest) (*pb.SeckillVoucherResponse, error) {
	userId := ctxdata.GetUidFromCtx(l.ctx)
	voucherId := req.VoucherId

	//orderId := l.svcCtx.RedisClient.Incr(l.ctx, "order:id").Val() // 使用 Redis ID 生成器
	redisIdWorker := tool.NewRedisIdWorker(redis.NewClient(&redis.Options{
		Addr:     l.svcCtx.Config.Cache[0].Host,
		Password: l.svcCtx.Config.Cache[0].Pass,
	}))
	orderId, _ := redisIdWorker.NextID("order")
	// 执行 Lua 脚本
	script, err := os.ReadFile("./deploy/script/seckill.lua")
	if err != nil {
		return nil, fmt.Errorf("读取脚本文件失败: %v", err)
	}

	result, err := l.svcCtx.RedisClient.Eval(l.ctx,
		string(script), nil,
		strconv.FormatInt(voucherId, 10),
		strconv.FormatInt(userId, 10),
		strconv.FormatInt(orderId, 10)).Result()

	if err != nil {
		l.Logger.Error("执行 Lua 脚本失败: ", err)
		return &pb.SeckillVoucherResponse{
			Code:    500,
			Message: "内部错误",
		}, nil
	}

	// 处理 Lua 脚本返回结果
	switch res := result.(type) {
	case int64:
		if res != 0 {
			return &pb.SeckillVoucherResponse{
				Code:    res,
				Message: getMessageByCode(res),
			}, nil
		}
	default:
		return &pb.SeckillVoucherResponse{
			Code:    500,
			Message: "内部错误",
		}, nil
	}

	return &pb.SeckillVoucherResponse{
		Code:    200,
		Message: "下单成功",
		OrderId: orderId,
	}, nil
}

// getMessageByCode 根据返回码获取消息
func getMessageByCode(code int64) string {
	switch code {
	case 1:
		return "库存不足"
	case 2:
		return "不能重复下单"
	default:
		return "未知错误"
	}
}
