package homestay

import (
	"context"
	"golodge/app/travel/cmd/api/internal/svc"
	"golodge/app/travel/cmd/api/internal/types"
	"golodge/app/travel/cmd/rpc/travel"
	"golodge/common/ctxdata"

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
