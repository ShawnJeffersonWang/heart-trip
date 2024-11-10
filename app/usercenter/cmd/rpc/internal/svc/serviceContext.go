package svc

import (
	"github.com/go-redis/redis/v8"
	"golodge/app/usercenter/cmd/rpc/internal/config"
	"golodge/app/usercenter/model"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config config.Config

	DB          *gorm.DB
	RedisClient *redis.Client

	UserModel     model.UserModel
	UserAuthModel model.UserAuthModel
}

func NewServiceContext(c config.Config) *ServiceContext {

	sqlConn := sqlx.NewMysql(c.DB.DataSource)
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
	return &ServiceContext{
		Config: c,
		//RedisClient: redis.New(c.Redis.Host, func(r *redis.Redis) {
		//	r.Type = c.Redis.Type
		//	r.Pass = c.Redis.Pass
		//}),
		DB:          db,
		RedisClient: redisClient,

		UserAuthModel: model.NewUserAuthModel(sqlConn, c.Cache),
		UserModel:     model.NewUserModel(sqlConn, c.Cache),
	}
}
