package svc

import (
	"golodge/app/travel/cmd/rpc/travel"
	"golodge/app/usercenter/cmd/rpc/usercenter"
	"golodge/app/websocket/cmd/api/internal/config"
	"golodge/app/websocket/cmd/api/internal/logic/ws"

	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config                config.Config
	Hub                   *ws.Hub
	UsercenterRpc         usercenter.Usercenter
	TravelRpc             travel.Travel
	SetUidToCtxMiddleware rest.Middleware
}

func NewServiceContext(c config.Config) *ServiceContext {
	hub := ws.NewHub()
	go hub.Run()

	return &ServiceContext{
		Config:        c,
		Hub:           hub,
		UsercenterRpc: usercenter.NewUsercenter(zrpc.MustNewClient(c.UsercenterRpcConf)),
		TravelRpc:     travel.NewTravel(zrpc.MustNewClient(c.TravelRpcConf)),
	}
}
