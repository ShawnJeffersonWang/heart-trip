package svc

import (
	"github.com/go-redis/redis/v8"
	"github.com/zeromicro/go-queue/kq"
	"golodge/app/travel/cmd/rpc/internal/config"
	"golodge/app/travel/model"
	"golodge/app/usercenter/cmd/rpc/usercenter"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config      config.Config
	DB          *gorm.DB
	RedisClient *redis.Client
	usercenter.Usercenter

	KqPusherClient *kq.Pusher

	HomestayModel         model.HomestayModel
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
		Addr:     c.Redis.Host,
		Password: c.Redis.Pass,
	})

	sqlConn := sqlx.NewMysql(c.DB.DataSource)

	return &ServiceContext{
		Config:         c,
		DB:             db,
		RedisClient:    redisClient,
		KqPusherClient: kq.NewPusher(c.KqPusherConf.Brokers, c.KqPusherConf.Topic),
		// 我了个骚刚bug: 没有初始化HomestayModel导致，travel模块的RPC调不动model
		HomestayModel:         model.NewHomestayModel(sqlConn, c.Cache),
		HomestayActivityModel: model.NewHomestayActivityModel(sqlConn, c.Cache),
		GuessModel:            model.NewGuessModel(sqlConn, c.Cache),
		UserHomestayModel:     model.NewUserHomestayModel(sqlConn, c.Cache),
		HistoryModel:          model.NewHistoryModel(sqlConn, c.Cache),
		UserHistoryModel:      model.NewUserHistoryModel(sqlConn, c.Cache),
		HomestayCommentModel:  model.NewHomestayCommentModel(sqlConn, c.Cache),
	}
}
