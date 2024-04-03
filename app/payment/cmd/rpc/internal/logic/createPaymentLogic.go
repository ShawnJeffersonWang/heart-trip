package logic

import (
	"context"
	"golodge/app/payment/cmd/rpc/internal/svc"
	"golodge/app/payment/cmd/rpc/pb"
	"golodge/app/payment/model"
	"golodge/common/xerr"
	"time"

	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
)

type CreatePaymentLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreatePaymentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreatePaymentLogic {
	return &CreatePaymentLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// CreatePayment create wechat pay prepayorder.
func (l *CreatePaymentLogic) CreatePayment(in *pb.CreatePaymentReq) (*pb.CreatePaymentResp, error) {

	data := new(model.ThirdPayment)
	//data.Sn = uniqueid.GenSn(uniqueid.SN_PREFIX_THIRD_PAYMENT)
	data.Sn = in.OrderSn
	data.UserId = in.UserId
	data.PayMode = in.PayModel
	data.PayTotal = in.PayTotal
	data.OrderSn = in.OrderSn
	data.ServiceType = model.ThirdPaymentServiceTypeHomestayOrder
	// bug: data.PayTime没有初始化
	data.PayTime = time.Now()

	_, err := l.svcCtx.ThirdPaymentModel.Insert(l.ctx, nil, data)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DB_ERROR), "create wechat pay prepayorder db insert fail , err:%v ,data : %+v  ", err, data)
	}

	return &pb.CreatePaymentResp{
		Sn: data.Sn,
	}, nil
}
