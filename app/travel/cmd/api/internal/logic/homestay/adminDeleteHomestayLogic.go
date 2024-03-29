package homestay

import (
	"context"
	"golodge/app/travel/cmd/api/internal/svc"
	"golodge/app/travel/cmd/api/internal/types"
	"golodge/app/travel/cmd/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type AdminDeleteHomestayLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAdminDeleteHomestayLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminDeleteHomestayLogic {
	return &AdminDeleteHomestayLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AdminDeleteHomestayLogic) AdminDeleteHomestay(req *types.AdminDeleteHomestayReq) (*types.AdminDeleteHomestayResp, error) {
	// todo: add your logic here and delete this line
	_, err := l.svcCtx.TravelRpc.AdminDeleteHomestay(l.ctx, &pb.AdminDeleteHomestayReq{
		HomestayId: req.HomestayId,
	})
	if err != nil {
		return nil, err
	}

	return &types.AdminDeleteHomestayResp{
		Success: true,
	}, nil
}
