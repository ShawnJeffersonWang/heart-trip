package logic

import (
	"context"
	"fmt"

	"heart-trip/app/order/cmd/rpc/internal/svc"
	"heart-trip/app/order/cmd/rpc/pb"
	"heart-trip/app/order/model"
	"heart-trip/common/xerr"

	"github.com/jinzhu/copier"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
)

type HomestayOrderDetailLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewHomestayOrderDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *HomestayOrderDetailLogic {
	return &HomestayOrderDetailLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *HomestayOrderDetailLogic) HomestayOrderDetail(in *pb.HomestayOrderDetailReq) (*pb.HomestayOrderDetailResp, error) {

	homestayOrder, err := l.svcCtx.HomestayOrderModel.FindOneBySn(l.ctx, in.Sn)
	if err != nil && err != model.ErrNotFound {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DB_ERROR), "HomestayOrderModel  FindOneBySn db err : %v , sn : %s", err, in.Sn)
	}

	var resp pb.HomestayOrder
	fmt.Println("homestayOrder: ", homestayOrder)
	if homestayOrder != nil {
		_ = copier.Copy(&resp, homestayOrder)

		resp.CreateTime = homestayOrder.CreateTime.Unix()
		resp.LiveStartDate = homestayOrder.LiveStartDate.Unix()
		resp.LiveEndDate = homestayOrder.LiveEndDate.Unix()

	}

	return &pb.HomestayOrderDetailResp{
		HomestayOrder: &resp,
	}, nil
}
