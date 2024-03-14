package svc

import (
	"homestay/app/order/cmd/api/internal/config"
	"homestay/app/order/cmd/rpc/order"
	"homestay/app/payment/cmd/rpc/payment"
	"homestay/app/travel/cmd/rpc/travel"

	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config config.Config

	OrderRpc   order.Order
	PaymentRpc payment.Payment
	TravelRpc  travel.Travel
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config: c,

		OrderRpc:   order.NewOrder(zrpc.MustNewClient(c.OrderRpcConf)),
		PaymentRpc: payment.NewPayment(zrpc.MustNewClient(c.PaymentRpcConf)),
		TravelRpc:  travel.NewTravel(zrpc.MustNewClient(c.TravelRpcConf)),
	}
}
