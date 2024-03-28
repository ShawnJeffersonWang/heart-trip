package homestay

import (
	"context"
	"golodge/app/travel/cmd/api/internal/svc"
	"golodge/app/travel/cmd/api/internal/types"
	"golodge/app/travel/cmd/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteHomestayLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteHomestayLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteHomestayLogic {
	return &DeleteHomestayLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteHomestayLogic) DeleteHomestay(req *types.DeleteHomestayReq) (*types.DeleteHomestayResp, error) {
	// todo: add your logic here and delete this line
	//userId := ctxdata.GetUidFromCtx(l.ctx)
	_, err := l.svcCtx.TravelRpc.DeleteHomestay(l.ctx, &pb.DeleteHomestayReq{
		//UserId:     userId,
		HomestayId: req.HomestayId,
	})
	if err != nil {
		return nil, err
	}

	return &types.DeleteHomestayResp{
		Success: true,
	}, nil
}
