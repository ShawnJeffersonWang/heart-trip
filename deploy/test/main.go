package main

import (
	"context"
	"encoding/json"
	"fmt"
	"golodge/common/globalkey"
	"golodge/common/tool"
	"log"
	"math/rand"
	"strconv"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
)

// 定义全局上下文
var ctx = context.Background()

// 常量
const (
	CacheShopKey      = "cache:shop:"
	HyperLogLogKey    = "hl2"
	LogicalExpireTime = 10 * time.Second
)

// Shop 模拟的 Shop 结构体
type Shop struct {
	ID     int64
	TypeID int64
	X      float64 // 经度
	Y      float64 // 纬度
	// 其他字段...
}

// CacheClient 模拟缓存客户端
type CacheClient struct {
	rdb *redis.Client
}

func NewCacheClient(rdb *redis.Client) *CacheClient {
	return &CacheClient{rdb: rdb}
}

func (c *CacheClient) SetWithLogicalExpire(key string, value interface{}, duration time.Duration) error {
	// 将 value 序列化为 JSON
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("序列化值失败: %v", err)
	}
	// 模拟设置带逻辑过期的缓存
	return c.rdb.Set(ctx, key, data, duration).Err()
}

// ShopServiceImpl 模拟商铺服务
type ShopServiceImpl struct {
	shops map[int64]*Shop
}

func NewShopService() *ShopServiceImpl {
	// 初始化一些模拟数据
	shops := make(map[int64]*Shop)
	for i := 1; i <= 10; i++ {
		shops[int64(i)] = &Shop{
			ID:     int64(i),
			TypeID: int64(rand.Intn(5) + 1),
			X:      rand.Float64()*360 - 180, // 经度范围：-180 到 +180
			Y:      rand.Float64()*180 - 90,  // 纬度范围：-90 到 +90
		}
	}
	return &ShopServiceImpl{shops: shops}
}

func (s *ShopServiceImpl) GetByID(id int64) (*Shop, error) {
	shop, exists := s.shops[id]
	if !exists {
		return nil, fmt.Errorf("shop with id %d not found", id)
	}
	return shop, nil
}

func (s *ShopServiceImpl) List() ([]*Shop, error) {
	list := make([]*Shop, 0, len(s.shops))
	for _, shop := range s.shops {
		list = append(list, shop)
	}
	return list, nil
}

// MainTest 包含所有测试方法
func MainTest() {
	// 初始化 Redis 客户端
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "shawn1234",
	})

	// 初始化依赖
	cacheClient := NewCacheClient(rdb)
	shopService := NewShopService()
	redisIdWorker := tool.NewRedisIdWorker(redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "shawn1234",
	}))

	// 执行各个测试
	testIdWorker(redisIdWorker)
	testSaveShop(cacheClient, shopService)
	loadShopData(rdb, shopService)
	testHyperLogLog(rdb)
}

// TestIdWorker 测试 ID 生成器
func testIdWorker(idWorker *tool.RedisIdWorker) {
	var wg sync.WaitGroup
	wg.Add(300)

	es := make(chan struct{}, 500) // 固定大小的协程池

	start := time.Now()

	for i := 0; i < 300; i++ {
		es <- struct{}{}
		go func() {
			defer wg.Done()
			for j := 0; j < 100; j++ {
				id, _ := idWorker.NextID("order")
				log.Printf("id = %d\n", id)
			}
			<-es
		}()
	}

	wg.Wait()
	elapsed := time.Since(start)
	log.Printf("time = %v\n", elapsed)
}

// TestSaveShop 测试保存店铺到缓存
func testSaveShop(cacheClient *CacheClient, shopService *ShopServiceImpl) {
	shop, err := shopService.GetByID(1)
	if err != nil {
		log.Printf("获取店铺失败: %v", err)
	}

	key := CacheShopKey + strconv.FormatInt(shop.ID, 10)
	err = cacheClient.SetWithLogicalExpire(key, shop, LogicalExpireTime)
	if err != nil {
		log.Printf("设置缓存失败: %v", err)
	}

	log.Println("店铺已缓存")
}

// LoadShopData 加载店铺数据到 Redis GEO
func loadShopData(rdb *redis.Client, shopService *ShopServiceImpl) {
	shops, err := shopService.List()
	if err != nil {
		log.Printf("获取店铺列表失败: %v", err)
	}

	// 按 TypeID 分组
	typeGroup := make(map[int64][]*Shop)
	for _, shop := range shops {
		typeGroup[shop.TypeID] = append(typeGroup[shop.TypeID], shop)
	}

	// 分批写入 Redis GEOADD
	for typeID, shops := range typeGroup {
		//key := ShopGeoKey + strconv.FormatInt(typeID, 10)
		key := fmt.Sprintf(globalkey.ShopGeoKey, typeID)
		var geoArgs []*redis.GeoLocation
		for _, shop := range shops {
			geoArgs = append(geoArgs, &redis.GeoLocation{
				Name:      strconv.FormatInt(shop.ID, 10),
				Longitude: shop.X,
				Latitude:  shop.Y,
			})
		}
		// 批处理，防止每一个店铺发一个请求
		err := rdb.GeoAdd(ctx, key, geoArgs...).Err()
		if err != nil {
			log.Printf("GEOADD 失败: %v", err)
		} else {
			log.Printf("已将 %d 个店铺添加到 %s\n", len(shops), key)
		}
	}
}

// TestHyperLogLog 测试 HyperLogLog
func testHyperLogLog(rdb *redis.Client) {
	values := make([]interface{}, 1000)
	for i := 0; i < 1000000; i++ {
		j := i % 1000
		values[j] = fmt.Sprintf("user_%d", i)
		if j == 999 {
			// 发送到 Redis
			err := rdb.PFAdd(ctx, HyperLogLogKey, values...).Err()
			if err != nil {
				log.Printf("PFADD 失败: %v", err)
			}
		}
	}

	// 统计数量
	count, err := rdb.PFCount(ctx, HyperLogLogKey).Result()
	if err != nil {
		log.Printf("PFCOUNT 失败: %v", err)
	}

	log.Printf("count = %d\n", count)
}

func main() {
	// 运行测试
	MainTest()
}
