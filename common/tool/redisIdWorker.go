package tool

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"time"
)

// RedisIdWorker 用于生成唯一ID
type RedisIdWorker struct {
	redisClient *redis.Client
	beginUnix   int64
	countBits   uint
	ctx         context.Context
}

// 常量定义
const (
	BeginTimestamp = 1640995200 // 2022-01-01 00:00:00 UTC
	CountBits      = 32
)

// NewRedisIdWorker 创建一个新的RedisIdWorker实例
func NewRedisIdWorker(redisClient *redis.Client) *RedisIdWorker {
	return &RedisIdWorker{
		redisClient: redisClient,
		beginUnix:   BeginTimestamp,
		countBits:   CountBits,
		ctx:         context.Background(),
	}
}

// NextID 生成下一个唯一ID
func (w *RedisIdWorker) NextID(keyPrefix string) (int64, error) {
	// 1. 生成时间戳
	now := time.Now().UTC()
	nowUnix := now.Unix()
	timestamp := nowUnix - w.beginUnix

	// 2. 生成序列号
	// 2.1 获取当前日期，精确到天
	date := now.Format("2006:01:02")
	// 2.2 自增长
	counterKey := fmt.Sprintf("icr:%s:%s", keyPrefix, date)
	count, err := w.redisClient.Incr(w.ctx, counterKey).Result()
	if err != nil {
		return 0, fmt.Errorf("failed to increment Redis key %s: %v", counterKey, err)
	}

	// 3. 拼接并返回
	id := (timestamp << w.countBits) | count
	return id, nil
}
