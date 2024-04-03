package thirdPayment

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
	"golodge/app/order/cmd/rpc/order"
	"golodge/app/payment/cmd/api/internal/svc"
	"golodge/app/payment/cmd/api/internal/types"
	"golodge/app/payment/cmd/rpc/payment"
	"golodge/app/payment/cmd/rpc/pb"
	"golodge/app/payment/model"
	"golodge/common/ctxdata"
	"golodge/common/xerr"
)

type FakePayCallbackLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

type NotifyTransaction struct {
	OutTradeNo     string
	PayerTotal     int64
	TradeState     string
	TransactionId  string
	TradeType      string
	TradeStateDesc string
}

func NewFakePayCallbackLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FakePayCallbackLogic {
	return &FakePayCallbackLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FakePayCallbackLogic) FakePayCallback(req *types.FakePayCallbackReq) (*types.FakePayCallbackResp, error) {
	// todo: add your logic here and delete this line
	var totalPrice int64
	var description string
	switch req.ServiceType {
	case model.ThirdPaymentServiceTypeHomestayOrder:
		homestayTotalPrice, homestayDescription, err := l.getPayHomestayPriceDescription(req.OrderSn)
		if err != nil {
			return nil, errors.New("getPayHomestayPriceDescription err")
		}
		totalPrice = homestayTotalPrice
		description = homestayDescription
	default:
		return nil, errors.New("Payment for this business type is not supported")
	}
	fmt.Println(totalPrice, "===", description)
	err := l.createWxPrePayOrder(req.ServiceType, req.OrderSn, totalPrice, description)
	if err != nil {
		return nil, err
	}

	transaction := new(NotifyTransaction)
	transaction.PayerTotal = req.PayerTotal
	transaction.OutTradeNo = req.OrderSn
	transaction.TradeState = "SUCCESS"
	returnCode := "SUCCESS"
	err = l.verifyAndUpdateState(transaction)
	if err != nil {
		returnCode = "FAIL"
	}
	return &types.FakePayCallbackResp{
		ReturnCode: returnCode,
	}, nil
}

func (l *FakePayCallbackLogic) verifyAndUpdateState(notifyTransaction *NotifyTransaction) error {
	paymentResp, err := l.svcCtx.PaymentRpc.GetPaymentBySn(l.ctx, &pb.GetPaymentBySnReq{
		Sn: notifyTransaction.OutTradeNo,
	})
	fmt.Println("paymentResp: ", paymentResp)
	if err != nil || paymentResp.PaymentDetail.Id == 0 {
		return errors.Wrapf(ErrWxPayCallbackError, "Failed to get payment flow record err:%v ,notifyTrasaction:%+v ", err, notifyTransaction)
	}
	notifyPayTotal := notifyTransaction.PayerTotal
	fmt.Println("notifyPayTotal: ", notifyTransaction.PayerTotal)
	fmt.Println("PayerTotal: ", notifyPayTotal)
	if notifyPayTotal != paymentResp.PaymentDetail.PayTotal {
		return errors.Wrapf(ErrWxPayCallbackError, "Order amount exception  notifyPayTotal:%v , notifyTrasaction:%v ", notifyPayTotal, notifyTransaction)
	}
	payStatus := l.getPayStatusByWXPayTradeState(notifyTransaction.TradeState)
	fmt.Println("payStatus: ", payStatus)
	if payStatus == model.ThirdPaymentPayTradeStateSuccess {
		if paymentResp.PaymentDetail.PayStatus != model.ThirdPaymentPayTradeStateWait {
			return nil
		}
		fmt.Println("-------------------------")
		if _, err := l.svcCtx.PaymentRpc.UpdateTradeState(l.ctx, &payment.UpdateTradeStateReq{
			Sn:             notifyTransaction.OutTradeNo,
			TradeState:     notifyTransaction.TradeState,
			TransactionId:  notifyTransaction.TransactionId,
			TradeType:      notifyTransaction.TradeType,
			TradeStateDesc: notifyTransaction.TradeStateDesc,
			PayStatus:      l.getPayStatusByWXPayTradeState(notifyTransaction.TradeState),
		}); err != nil {
			return errors.Wrapf(ErrWxPayCallbackError, "更新流水状态失败: err: %v, notifyTransaction: %v", err, notifyTransaction)
		}
	} else if payStatus == model.ThirdPaymentPayTradeStateWait {
		// Refund notification
	}
	return nil
}

func (l *FakePayCallbackLogic) getPayStatusByWXPayTradeState(wxPayTradeState string) int64 {
	switch wxPayTradeState {
	case SUCCESS:
		return model.ThirdPaymentPayTradeStateSuccess
	case USERPAYING:
		return model.ThirdPaymentPayTradeStateWait
	case REFUND:
		return model.ThirdPaymentPayTradeStateWait
	default:
		return model.ThirdPaymentPayTradeStateFAIL
	}
}

func (l *FakePayCallbackLogic) createWxPrePayOrder(serviceType, orderSn string, totalPrice int64, description string) error {
	userId := ctxdata.GetUidFromCtx(l.ctx)

	// 持久化
	createPaymentResp, err := l.svcCtx.PaymentRpc.CreatePayment(l.ctx, &payment.CreatePaymentReq{
		UserId:      userId,
		PayModel:    model.ThirdPaymentPayModelWechatPay,
		PayTotal:    totalPrice,
		OrderSn:     orderSn,
		ServiceType: serviceType,
	})
	if err != nil || createPaymentResp.Sn == "" {
		return errors.WithMessagef(err, "create local third payment record fail: err: %v, userId: %d, totalPrice: %d, orderSn: %s",
			err, userId, totalPrice, orderSn)
	}
	return nil
}

func (l *FakePayCallbackLogic) getPayHomestayPriceDescription(orderSn string) (int64, string, error) {
	description := "homestay pay"

	resp, err := l.svcCtx.OrderRpc.HomestayOrderDetail(l.ctx, &order.HomestayOrderDetailReq{
		Sn: orderSn,
	})
	if err != nil {
		return 0, description, errors.Wrapf(xerr.NewErrMsg("order no exists"), "wechat payment order does not exist")
	}
	return resp.HomestayOrder.OrderTotalPrice, description, nil
}
