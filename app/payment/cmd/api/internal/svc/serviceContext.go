package svc

import (
	"heart-trip/app/order/cmd/rpc/order"
	"heart-trip/app/payment/cmd/api/internal/config"
	"heart-trip/app/payment/cmd/rpc/payment"
	"heart-trip/app/usercenter/cmd/rpc/usercenter"

	"github.com/wechatpay-apiv3/wechatpay-go/core"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config config.Config

	WxPayClient *core.Client

	PaymentRpc    payment.Payment
	OrderRpc      order.Order
	UsercenterRpc usercenter.Usercenter
}

func NewServiceContext(c config.Config) *ServiceContext {

	return &ServiceContext{
		Config: c,

		PaymentRpc:    payment.NewPayment(zrpc.MustNewClient(c.PaymentRpcConf)),
		OrderRpc:      order.NewOrder(zrpc.MustNewClient(c.OrderRpcConf)),
		UsercenterRpc: usercenter.NewUsercenter(zrpc.MustNewClient(c.UsercenterRpcConf)),
	}
}
