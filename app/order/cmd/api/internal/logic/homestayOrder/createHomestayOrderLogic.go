package homestayOrder

import (
	"context"
	"heart-trip/app/travel/cmd/rpc/pb"
	"heart-trip/common/ctxdata"

	"heart-trip/app/order/cmd/api/internal/svc"
	"heart-trip/app/order/cmd/api/internal/types"
	"heart-trip/app/order/cmd/rpc/order"
	"heart-trip/common/xerr"

	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
)

type CreateHomestayOrderLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateHomestayOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) CreateHomestayOrderLogic {
	return CreateHomestayOrderLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// create order
func (l *CreateHomestayOrderLogic) CreateHomestayOrder(req types.CreateHomestayOrderReq) (*types.CreateHomestayOrderResp, error) {

	homestayResp, err := l.svcCtx.TravelRpc.HomestayDetail(l.ctx, &pb.HomestayDetailReq{
		HomestayId: req.HomestayId,
	})
	if err != nil {
		return nil, err
	}
	//if homestayResp.Homestay == nil || homestayResp.Homestay.Id == 0 {
	//	return nil, errors.Wrapf(xerr.NewErrMsg("homestay no exists"), "CreateHomestayOrder homestay no exists id : %d", req.HomestayId)
	//}
	if homestayResp == nil || homestayResp.Id == 0 {
		return nil, errors.Wrapf(xerr.NewErrMsg("homestay no exists"), "CreateHomestayOrder homestay no exists id : %d", req.HomestayId)
	}

	userId := ctxdata.GetUidFromCtx(l.ctx)

	resp, err := l.svcCtx.OrderRpc.CreateHomestayOrder(l.ctx, &order.CreateHomestayOrderReq{
		HomestayId:    req.HomestayId,
		LiveStartTime: req.LiveStartTime,
		LiveEndTime:   req.LiveEndTime,
		UserId:        userId,
		Remark:        req.Remark,
	})
	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrMsg("create homestay order fail"), "create homestay order rpc CreateHomestayOrder fail req: %+v , err : %v ", req, err)
	}

	return &types.CreateHomestayOrderResp{
		OrderSn: resp.Sn,
	}, nil
}
