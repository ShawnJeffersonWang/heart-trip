package svc

import (
	"homestay/app/travel/cmd/rpc/travel"
	"homestay/app/usercenter/cmd/api/internal/config"
	"homestay/app/usercenter/cmd/rpc/usercenter"

	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config                config.Config
	UsercenterRpc         usercenter.Usercenter
	TravelRpc             travel.Travel
	SetUidToCtxMiddleware rest.Middleware
}

func NewServiceContext(c config.Config) *ServiceContext {

	return &ServiceContext{
		Config:        c,
		UsercenterRpc: usercenter.NewUsercenter(zrpc.MustNewClient(c.UsercenterRpcConf)),
		TravelRpc:     travel.NewTravel(zrpc.MustNewClient(c.TravelRpcConf)),
	}
}
