package svc

import (
	"homestay/app/order/cmd/mq/internal/config"
	"homestay/app/order/cmd/rpc/order"
	"homestay/app/usercenter/cmd/rpc/usercenter"

	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config config.Config

	OrderRpc      order.Order
	UsercenterRpc usercenter.Usercenter
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config: c,

		OrderRpc:      order.NewOrder(zrpc.MustNewClient(c.OrderRpcConf)),
		UsercenterRpc: usercenter.NewUsercenter(zrpc.MustNewClient(c.UsercenterRpcConf)),
	}
}
