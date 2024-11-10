package svc

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v8"
	"github.com/zeromicro/go-queue/kq"
	"github.com/zeromicro/go-zero/core/logx"
	"golodge/app/travel/cmd/rpc/internal/config"
	"golodge/app/travel/model"
	"golodge/app/usercenter/cmd/rpc/usercenter"
	"golodge/common/tool"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"os"
	"path/filepath"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config        config.Config
	DB            *gorm.DB
	RedisClient   *redis.Client
	CacheClient   *tool.CacheClient
	RedSync       *redsync.Redsync
	RedisIdWorker *tool.RedisIdWorker
	//LuaScripts  map[string]string // 存储脚本的SHA1

	usercenter.Usercenter

	KqPusherClient *kq.Pusher

	HomestayModel         model.HomestayModel
	ShopModel             model.ShopModel
	HomestayActivityModel model.HomestayActivityModel
	GuessModel            model.GuessModel
	UserHomestayModel     model.UserHomestayModel
	HistoryModel          model.HistoryModel
	UserHistoryModel      model.UserHistoryModel
	HomestayCommentModel  model.HomestayCommentModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	// 初始化数据库
	db, err := gorm.Open(mysql.Open(c.DB.DataSource), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info), // 根据需要调整日志级别
	})
	if err != nil {
		panic(err)
	}
	// 初始化 Redis 客户端
	redisClient := redis.NewClient(&redis.Options{
		Addr:     c.Cache[0].Host,
		Password: c.Cache[0].Pass,
	})
	cacheClient := tool.NewCacheClient(redisClient, logx.WithContext(context.Background()))

	// 初始化 RedSync
	pools := goredis.NewPool(redisClient)
	r := redsync.New(pools)

	redisIdWorker := tool.NewRedisIdWorker(redisClient)

	// 加载 Lua 脚本
	//luaScripts, err := loadLuaScripts(c.LuaScripts, redisClient)
	//if err != nil {
	//	log.Fatalf("加载 Lua 脚本失败: %v", err)
	//}

	sqlConn := sqlx.NewMysql(c.DB.DataSource)
	serviceContext := &ServiceContext{
		Config:        c,
		DB:            db,
		RedisClient:   redisClient,
		CacheClient:   cacheClient,
		RedSync:       r,
		RedisIdWorker: redisIdWorker,
		//LuaScripts:     luaScripts,
		KqPusherClient: kq.NewPusher(c.KqPusherConf.Brokers, c.KqPusherConf.Topic),
		// 我了个骚刚bug: 没有初始化HomestayModel导致，travel模块的RPC调不动model
		HomestayModel:         model.NewHomestayModel(sqlConn, c.Cache),
		ShopModel:             model.NewShopModel(db),
		HomestayActivityModel: model.NewHomestayActivityModel(sqlConn, c.Cache),
		GuessModel:            model.NewGuessModel(sqlConn, c.Cache),
		UserHomestayModel:     model.NewUserHomestayModel(sqlConn, c.Cache),
		HistoryModel:          model.NewHistoryModel(sqlConn, c.Cache),
		UserHistoryModel:      model.NewUserHistoryModel(sqlConn, c.Cache),
		HomestayCommentModel:  model.NewHomestayCommentModel(sqlConn, c.Cache),
	}
	go NewVoucherOrderHandler(c, serviceContext, redisClient).Start()
	return serviceContext
}

// loadLuaScripts 读取并加载 Lua 脚本，返回脚本的 SHA1 哈希
func loadLuaScripts(luaScriptPaths map[string]string, redisClient *redis.Client) (map[string]string, error) {
	luaScripts := make(map[string]string)
	for name, path := range luaScriptPaths {
		absolutePath, err := filepath.Abs(path)
		if err != nil {
			return nil, fmt.Errorf("获取脚本路径失败: %v", err)
		}

		scriptContent, err := os.ReadFile(absolutePath)
		if err != nil {
			return nil, fmt.Errorf("读取脚本文件失败: %v", err)
		}

		// 加载脚本到 Redis 并获取 SHA1
		sha, err := redisClient.ScriptLoad(context.Background(), string(scriptContent)).Result()
		if err != nil {
			return nil, fmt.Errorf("加载脚本到 Redis 失败: %v", err)
		}

		luaScripts[name] = sha
	}
	return luaScripts, nil
}
