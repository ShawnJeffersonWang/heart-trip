package logic

import (
	"context"

	"golodge/app/travel/cmd/rpc/internal/svc"
	"golodge/app/travel/cmd/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type RemoveWishListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRemoveWishListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RemoveWishListLogic {
	return &RemoveWishListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *RemoveWishListLogic) RemoveWishList(in *pb.RemoveWishListReq) (*pb.RemoveWishListResp, error) {
	// todo: add your logic here and delete this line
	err := l.svcCtx.UserHomestayModel.Delete(l.ctx, in.UserId, in.HomestayId)
	if err != nil {
		return nil, err
	}

	return &pb.RemoveWishListResp{}, nil
}
