package homestay

import (
	"context"
	"golodge/app/travel/cmd/api/internal/svc"
	"golodge/app/travel/cmd/api/internal/types"
	"golodge/app/travel/cmd/rpc/travel"

	"github.com/zeromicro/go-zero/core/logx"
)

type SeckillVoucherLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// seckill voucher order
func NewSeckillVoucherLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SeckillVoucherLogic {
	return &SeckillVoucherLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SeckillVoucherLogic) SeckillVoucher(req *types.SeckillVoucherRequest) (resp *types.SeckillVoucherResponse, err error) {
	seckillVoucherResponse, err := l.svcCtx.TravelRpc.SeckillVoucher(l.ctx, &travel.SeckillVoucherRequest{
		VoucherId: req.VoucherId,
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
