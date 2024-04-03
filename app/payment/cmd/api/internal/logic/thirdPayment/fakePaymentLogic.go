package thirdPayment

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"github.com/wechatpay-apiv3/wechatpay-go/core"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments/jsapi"
	"github.com/wechatpay-apiv3/wechatpay-go/utils"
	order "golodge/app/order/cmd/rpc/order"
	"golodge/app/payment/cmd/api/internal/svc"
	"golodge/app/payment/cmd/api/internal/types"
	"golodge/app/payment/cmd/rpc/payment"
	"golodge/app/payment/model"
	"golodge/common/ctxdata"
	"golodge/common/xerr"
	"strconv"
	"time"

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
	wechatPrepayRsp, err := l.createWxPrePayOrder(req.ServiceType, req.OrderSn, totalPrice, description)
	if err != nil {
		return nil, err
	}
	return &types.FakePaymentResp{
		Appid:     l.svcCtx.Config.WxMiniConf.AppId,
		NonceStr:  *wechatPrepayRsp.NonceStr,
		PaySign:   *wechatPrepayRsp.PaySign,
		Package:   *wechatPrepayRsp.Package,
		Timestamp: *wechatPrepayRsp.TimeStamp,
		SignType:  *wechatPrepayRsp.SignType,
	}, nil
}

func (l *FakePaymentLogic) createWxPrePayOrder(serviceType, orderSn string, totalPrice int64, description string) (*jsapi.PrepayWithRequestPaymentResponse, error) {
	userId := ctxdata.GetUidFromCtx(l.ctx)
	//userResp, err := l.svcCtx.UsercenterRpc.GetUserAuthByUserId(l.ctx, &usercenter.GetUserAuthByUserIdReq{
	//	UserId:   userId,
	//	AuthType: usercenterModel.UserAuthTypeSmallWX,
	//})
	//if err != nil {
	//	return nil, errors.New("Get user wechat openid err")
	//}
	//if userResp.UserAuth == nil || userResp.UserAuth.Id == 0 {
	//	return nil, errors.WithMessage(err, "Get user wechat openid fail, Please pay before authorization by wechat")
	//}
	//openId := userResp.UserAuth.AuthKey

	// 持久化
	createPaymentResp, err := l.svcCtx.PaymentRpc.CreatePayment(l.ctx, &payment.CreatePaymentReq{
		UserId:      userId,
		PayModel:    model.ThirdPaymentPayModelWechatPay,
		PayTotal:    totalPrice,
		OrderSn:     orderSn,
		ServiceType: serviceType,
	})
	if err != nil || createPaymentResp.Sn == "" {
		return nil, errors.WithMessagef(err, "create local third payment record fail: err: %v, userId: %d, totalPrice: %d, orderSn: %s",
			err, userId, totalPrice, orderSn)
	}
	nonce, err := utils.GenerateNonce()
	if err != nil {
		return nil, fmt.Errorf("generate request for payment err: %s", err.Error())
	}
	resp := &jsapi.PrepayWithRequestPaymentResponse{
		PrepayId:  core.String("PrepayId"),
		Appid:     core.String(l.svcCtx.Config.WxMiniConf.AppId),
		TimeStamp: core.String(strconv.FormatInt(time.Now().Unix(), 10)),
		NonceStr:  core.String(nonce),
		Package:   core.String("Package"),
		SignType:  core.String("RSA"),
		PaySign:   core.String("PaySign"),
	}
	return resp, err
}

func (l *FakePaymentLogic) getPayHomestayPriceDescription(orderSn string) (int64, string, error) {
	description := "homestay pay"

	resp, err := l.svcCtx.OrderRpc.HomestayOrderDetail(l.ctx, &order.HomestayOrderDetailReq{
		Sn: orderSn,
	})
	if err != nil {
		return 0, description, errors.Wrapf(xerr.NewErrMsg("order no exists"), "wechat payment order does not exist")
	}
	return resp.HomestayOrder.OrderTotalPrice, description, nil
}
