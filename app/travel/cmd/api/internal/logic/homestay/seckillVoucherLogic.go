package homestay

import (
	"context"
	"heart-trip/app/travel/cmd/api/internal/svc"
	"heart-trip/app/travel/cmd/api/internal/types"
	"heart-trip/app/travel/cmd/rpc/travel"
	"heart-trip/common/ctxdata"

	"github.com/zeromicro/go-zero/core/logx"
)

type SeckillVoucherLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSeckillVoucherLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SeckillVoucherLogic {
	return &SeckillVoucherLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SeckillVoucherLogic) SeckillVoucher(req *types.SeckillVoucherRequest) (resp *types.SeckillVoucherResponse, err error) {
	userId := ctxdata.GetUidFromCtx(l.ctx)
	seckillVoucherResponse, err := l.svcCtx.TravelRpc.SeckillVoucher(l.ctx, &travel.SeckillVoucherRequest{
		VoucherId: req.VoucherId,
		UserId:    userId,
	})
	if err != nil {
		return nil, err
	}
	return &types.SeckillVoucherResponse{
		Code:    seckillVoucherResponse.Code,
		Message: seckillVoucherResponse.Message,
		OrderId: seckillVoucherResponse.OrderId,
	}, nil
}
