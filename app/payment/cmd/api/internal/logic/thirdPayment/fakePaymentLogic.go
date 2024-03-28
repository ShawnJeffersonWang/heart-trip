package thirdPayment

import (
	"context"
	"golodge/app/payment/cmd/rpc/pb"

	"golodge/app/payment/cmd/api/internal/svc"
	"golodge/app/payment/cmd/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type FakePaymentLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFakePaymentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FakePaymentLogic {
	return &FakePaymentLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FakePaymentLogic) FakePayment(req *types.FakePaymentReq) (resp *types.FakePaymentResp, err error) {
	// todo: add your logic here and delete this line
	l.svcCtx.PaymentRpc.UpdateTradeState(l.ctx, &pb.UpdateTradeStateReq{
		Sn:             req.OrderSn,
		TradeState:     "",
		TransactionId:  "",
		TradeType:      req.ServiceType,
		TradeStateDesc: "",
		PayStatus:      1,
		PayTime:        0,
	})
	return &types.FakePaymentResp{
		Appid:     "",
		NonceStr:  "",
		PaySign:   "",
		Package:   "",
		Timestamp: "",
		SignType:  "",
	}, nil
}
