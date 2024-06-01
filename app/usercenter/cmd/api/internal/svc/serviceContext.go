package svc

import (
	"golodge/app/travel/cmd/rpc/travel"
	"golodge/app/usercenter/cmd/api/internal/config"
	"golodge/app/usercenter/cmd/api/internal/logic/ws"
	"golodge/app/usercenter/cmd/rpc/usercenter"

	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config                config.Config
	Hub                   *ws.Hub
	UsercenterRpc         usercenter.Usercenter
	TravelRpc             travel.Travel
	SetUidToCtxMiddleware rest.Middleware
	Hua                   *ws.Hub
}

func NewServiceContext(c config.Config) *ServiceContext {
	hub := ws.NewHub()
	go hub.Run()

	return &ServiceContext{
		Config:        c,
		UsercenterRpc: usercenter.NewUsercenter(zrpc.MustNewClient(c.UsercenterRpcConf)),
		TravelRpc:     travel.NewTravel(zrpc.MustNewClient(c.TravelRpcConf)),
	}
}
