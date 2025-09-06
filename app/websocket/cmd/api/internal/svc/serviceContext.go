package svc

import (
	"heart-trip/app/travel/cmd/rpc/travel"
	"heart-trip/app/usercenter/cmd/rpc/usercenter"
	"heart-trip/app/websocket/cmd/api/internal/config"
	"heart-trip/app/websocket/model"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config                config.Config
	Hub                   *Hub
	UsercenterRpc         usercenter.Usercenter
	MessageModel          model.MessageModel
	TravelRpc             travel.Travel
	SetUidToCtxMiddleware rest.Middleware
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := sqlx.NewMysql(c.DB.DataSource)

	hub := NewHub()
	go hub.Run()

	return &ServiceContext{
		Config:        c,
		Hub:           hub,
		UsercenterRpc: usercenter.NewUsercenter(zrpc.MustNewClient(c.UsercenterRpcConf)),
		MessageModel:  model.NewMessageModel(conn),
		TravelRpc:     travel.NewTravel(zrpc.MustNewClient(c.TravelRpcConf)),
	}
}
