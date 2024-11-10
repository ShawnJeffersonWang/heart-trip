package tool

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/zeromicro/go-zero/core/logx"
	"golodge/common/globalkey"
	"time"
)

// CacheClient 封装 Redis 操作及缓存逻辑
type CacheClient struct {
	redisClient           *redis.Client
	cacheRebuildSemaphore chan struct{} // 限制缓存重建的并发数量
	logger                logx.Logger
}

// NewCacheClient 创建一个新的 CacheClient 实例
func NewCacheClient(redisClient *redis.Client, logger logx.Logger) *CacheClient {
	return &CacheClient{
		redisClient:           redisClient,
		cacheRebuildSemaphore: make(chan struct{}, 10), // 最大并发数为10
		logger:                logger,
	}
}

// Set 将数据序列化为 JSON 并设置到 Redis 中，指定过期时间
func (c *CacheClient) Set(ctx context.Context, key string, value interface{}, t time.Duration) error {
	jsonData, err := json.Marshal(value)
	if err != nil {
		c.logger.Errorf("Set JSON Marshal error: %v", err)
		return err
	}
	return c.redisClient.Set(ctx, key, jsonData, t).Err()
}

// RedisData 用于封装缓存数据和逻辑过期时间
type RedisData struct {
	Data       json.RawMessage `json:"data"`
	ExpireTime time.Time       `json:"expire_time"`
}

// SetWithLogicalExpire 将数据和逻辑过期时间封装后设置到 Redis 中，不设置 Redis 过期时间
func (c *CacheClient) SetWithLogicalExpire(ctx context.Context, key string, value interface{}, t time.Duration) error {
	redisData := RedisData{
		Data:       mustMarshalJSON(value),
		ExpireTime: time.Now().Add(t),
	}
	jsonData, err := json.Marshal(redisData)
	if err != nil {
		c.logger.Errorf("SetWithLogicalExpire JSON Marshal error: %v", err)
		return err
	}
	return c.redisClient.Set(ctx, key, jsonData, 0).Err() // 不设置 Redis 过期时间
}

func mustMarshalJSON(v interface{}) json.RawMessage {
	data, err := json.Marshal(v)
	if err != nil {
		panic(fmt.Sprintf("MustMarshalJSON failed: %v", err))
	}
	return data
}

// QueryWithPassThrough 从 Redis 查询缓存，如果未命中则回源数据库，并缓存结果
func (c *CacheClient) QueryWithPassThrough(
	ctx context.Context,
	keyPrefix string,
	id any,
	dbFallback func(any) (any, error),
	t time.Duration,
) (any, error) {
	var zero any
	key := fmt.Sprintf("%s%v", keyPrefix, id)

	// 1. 从 Redis 查询缓存
	jsonData, err := c.redisClient.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		// 缓存未命中，继续回源
	} else if err != nil {
		c.logger.Errorf("QueryWithPassThrough Redis Get error: %v", err)
		return zero, err
	} else {
		// 2. 判断缓存是否存在且不为空
		if len(jsonData) > 0 {
			var res any
			if err := json.Unmarshal([]byte(jsonData), &res); err == nil {
				return res, nil
			}
			c.logger.Errorf("QueryWithPassThrough JSON Unmarshal error: %v", err)
		} else {
			// 缓存存在但为空，表示查询结果为空
			return zero, nil
		}
	}

	// 3. 缓存未命中，回源数据库
	res, err := dbFallback(id)
	if err != nil {
		c.logger.Errorf("QueryWithPassThrough dbFallback error: %v", err)
		return zero, err
	}
	if res == nil {
		// 4. 数据库查询结果为空，缓存空值
		if err := c.redisClient.Set(ctx, key, "", globalkey.CacheNullTtl).Err(); err != nil {
			c.logger.Errorf("QueryWithPassThrough Redis Set empty value error for key %s: %v", key, err)
		}
		return zero, nil
	}

	// 5. 数据库查询结果存在，缓存数据
	if err := c.Set(ctx, key, res, t); err != nil {
		c.logger.Errorf("QueryWithPassThrough Redis Set error: %v", err)
	}
	return res, nil
}

// QueryWithLogicalExpire 从 Redis 查询缓存，并根据逻辑过期时间决定是否异步重建缓存
func (c *CacheClient) QueryWithLogicalExpire(
	ctx context.Context,
	keyPrefix string,
	id any,
	dbFallback func(any) (any, error),
	t time.Duration,
) (any, error) {
	var zero any
	key := fmt.Sprintf("%s%v", keyPrefix, id)

	// 1. 从 Redis 查询缓存
	jsonData, err := c.redisClient.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		// 缓存未命中，返回零值
		return zero, nil
	} else if err != nil {
		c.logger.Errorf("QueryWithLogicalExpire Redis Get error: %v", err)
		return zero, err
	}

	// 2. 解析 RedisData
	var redisData RedisData
	if err := json.Unmarshal([]byte(jsonData), &redisData); err != nil {
		c.logger.Errorf("QueryWithLogicalExpire JSON Unmarshal error: %v", err)
		return zero, err
	}

	// 3. 判断逻辑过期时间
	if redisData.ExpireTime.After(time.Now()) {
		// 未过期，返回数据
		var res any
		if err := json.Unmarshal(redisData.Data, &res); err != nil {
			c.logger.Errorf("QueryWithLogicalExpire JSON Unmarshal Data error: %v", err)
			return zero, err
		}
		return res, nil
	}

	// 4. 已过期，尝试获取锁
	lockKey := fmt.Sprintf("%s%v", globalkey.LockShopKey, id)
	gotLock, err := c.tryLock(ctx, lockKey)
	if err != nil {
		c.logger.Errorf("QueryWithLogicalExpire tryLock error: %v", err)
	}

	if gotLock {
		// 5. 获取锁成功，尝试重建缓存
		select {
		case c.cacheRebuildSemaphore <- struct{}{}:
			// 在独立的 goroutine 中重建缓存
			go func() {
				defer func() { <-c.cacheRebuildSemaphore }()
				defer c.unlock(ctx, lockKey)

				// 回源数据库
				newRes, err := dbFallback(id)
				if err != nil {
					c.logger.Errorf("QueryWithLogicalExpire dbFallback error: %v", err)
					return
				}
				if newRes == nil {
					// 数据库查询结果为空，不重建缓存
					return
				}
				// 更新缓存
				if err := c.SetWithLogicalExpire(ctx, key, newRes, t); err != nil {
					c.logger.Errorf("QueryWithLogicalExpire SetWithLogicalExpire error: %v", err)
				}
			}()
		default:
			// 缓存重建队列已满，跳过缓存重建
			c.logger.Errorf("QueryWithLogicalExpire cache rebuild executor is full")
		}
	}

	// 6. 返回过期的数据
	var res any
	if err := json.Unmarshal(redisData.Data, &res); err != nil {
		c.logger.Errorf("QueryWithLogicalExpire JSON Unmarshal Data error: %v", err)
		return zero, err
	}
	return res, nil
}

// QueryWithMutex 从 Redis 查询缓存，如果未命中则加锁回源并缓存，防止缓存击穿
func (c *CacheClient) QueryWithMutex(
	ctx context.Context,
	keyPrefix string,
	id any,
	dbFallback func(any) (any, error),
	t time.Duration,
) (any, error) {
	var zero any
	key := fmt.Sprintf("%s%v", keyPrefix, id)

	// 1. 从 Redis 查询缓存
	jsonData, err := c.redisClient.Get(ctx, key).Result()
	if !errors.Is(err, redis.Nil) && err != nil {
		c.logger.Errorf("QueryWithMutex Redis Get error: %v", err)
		return zero, err
	}

	if err == nil && len(jsonData) > 0 {
		// 缓存命中，返回数据
		var res any
		if err := json.Unmarshal([]byte(jsonData), &res); err == nil {
			return res, nil
		}
		c.logger.Errorf("QueryWithMutex JSON Unmarshal error: %v", err)
	}

	if err == nil && jsonData == "" {
		// 缓存存在但为空，表示查询结果为空
		return zero, nil
	}

	// 2. 尝试获取锁
	lockKey := fmt.Sprintf("%s%v", globalkey.LockShopKey, id)
	gotLock, err := c.tryLock(ctx, lockKey)
	if err != nil {
		c.logger.Errorf("QueryWithMutex tryLock error: %v", err)
		return zero, err
	}

	if !gotLock {
		// 获取锁失败，休眠并重试
		time.Sleep(50 * time.Millisecond)
		return c.QueryWithMutex(ctx, keyPrefix, id, dbFallback, t)
	}

	// 3. 获取锁成功，确保释放锁
	defer c.unlock(ctx, lockKey)

	// 4. 回源数据库
	res, err := dbFallback(id)
	if err != nil {
		c.logger.Errorf("QueryWithMutex dbFallback error: %v", err)
		return zero, err
	}
	if res == nil {
		// 5. 数据库查询结果为空，缓存空值
		if err := c.redisClient.Set(ctx, key, "", globalkey.CacheNullTtl).Err(); err != nil {
			c.logger.Errorf("QueryWithMutex Redis Set empty value error: %v", err)
		}
		return zero, nil
	}

	// 6. 数据库查询结果存在，缓存数据
	if err := c.Set(ctx, key, res, t); err != nil {
		c.logger.Errorf("QueryWithMutex Redis Set error: %v", err)
	}

	return res, nil
}

// tryLock 尝试获取锁，设置锁的过期时间为 10 秒
func (c *CacheClient) tryLock(ctx context.Context, key string) (bool, error) {
	return c.redisClient.SetNX(ctx, key, "1", 10*time.Second).Result()
}

// unlock 释放锁
func (c *CacheClient) unlock(ctx context.Context, key string) error {
	return c.redisClient.Del(ctx, key).Err()
}
