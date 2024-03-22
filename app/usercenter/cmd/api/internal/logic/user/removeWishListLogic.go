package user

import (
	"context"
	"golodge/app/travel/cmd/rpc/travel"
	"golodge/common/ctxdata"

	"golodge/app/usercenter/cmd/api/internal/svc"
	"golodge/app/usercenter/cmd/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type RemoveWishListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRemoveWishListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RemoveWishListLogic {
	return &RemoveWishListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RemoveWishListLogic) RemoveWishList(req *types.RemoveWishListReq) (*types.RemoveWishListResp, error) {
	// todo: add your logic here and delete this line
	userId := ctxdata.GetUidFromCtx(l.ctx)

	_, err := l.svcCtx.TravelRpc.RemoveWishList(l.ctx, &travel.RemoveWishListReq{
		UserId:     userId,
		HomestayId: req.Id,
	})
	if err != nil {
		return nil, err
	}
	return &types.RemoveWishListResp{
		Success: true,
	}, nil
}
