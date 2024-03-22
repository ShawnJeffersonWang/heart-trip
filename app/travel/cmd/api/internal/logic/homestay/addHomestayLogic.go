package homestay

import (
	"context"
	"golodge/app/travel/cmd/api/internal/svc"
	"golodge/app/travel/cmd/api/internal/types"
	"golodge/app/travel/cmd/rpc/pb"
	"golodge/app/travel/cmd/rpc/travel"
	"golodge/common/ctxdata"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddHomestayLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAddHomestayLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddHomestayLogic {
	return &AddHomestayLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AddHomestayLogic) AddHomestay(req *types.AddHomestayReq) (*types.AddHomestayResp, error) {
	// todo: add your logic here and delete this line
	userId := ctxdata.GetUidFromCtx(l.ctx)

	_, err := l.svcCtx.TravelRpc.AddHomestay(l.ctx, &travel.AddHomestayReq{
		Homestay: &pb.Homestay{
			UserId:      userId,
			Title:       req.Homestay.Title,
			Cover:       req.Homestay.Cover,
			Intro:       req.Homestay.Intro,
			Location:    req.Homestay.Location,
			PriceBefore: req.Homestay.PriceBefore,
		},
	})
	if err != nil {
		return nil, err
	}

	return &types.AddHomestayResp{
		Success: true,
	}, nil
}
