package logic

import (
	"context"
	"github.com/zeromicro/go-zero/core/stores/sqlx"

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
	userHomestay, err := l.svcCtx.UserHomestayModel.FindOneByUserIdAndHomestayId(l.ctx, in.UserId, in.HomestayId)
	err = l.svcCtx.UserHistoryModel.Trans(l.ctx, func(ctx context.Context, session sqlx.Session) error {
		err := l.svcCtx.UserHomestayModel.DeleteSoft(ctx, session, userHomestay)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &pb.RemoveWishListResp{}, nil
}
