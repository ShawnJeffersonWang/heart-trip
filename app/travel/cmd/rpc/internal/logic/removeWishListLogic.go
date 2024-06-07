package logic

import (
	"context"
	"github.com/zeromicro/go-zero/core/logx"
	"golodge/app/travel/cmd/rpc/internal/svc"
	"golodge/app/travel/cmd/rpc/pb"
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
	//userHomestay, err := l.svcCtx.UserHomestayModel.FindOneByUserIdAndHomestayId(l.ctx, in.UserId, in.HomestayId)
	//err = l.svcCtx.UserHistoryModel.Trans(l.ctx, func(ctx context.Context, session sqlx.Session) error {
	//	err := l.svcCtx.UserHomestayModel.DeleteSoft(ctx, session, userHomestay)
	//	if err != nil {
	//		return err
	//	}
	//	return nil
	//})
	//if err != nil {
	//	return nil, err
	//}
	//return &pb.RemoveWishListResp{}, nil

	// 检查用户是否已经收藏了该民宿
	exists, err := l.svcCtx.UserHomestayModel.CheckIfExists(l.ctx, in.UserId, in.HomestayId)
	if err != nil {
		return nil, err
	}
	if !exists {
		return &pb.RemoveWishListResp{
			Success: false,
		}, nil
	}

	// 更新 del_state 字段以取消收藏
	err = l.svcCtx.UserHomestayModel.UpdateDelState(l.ctx, in.UserId, in.HomestayId, 1)
	if err != nil {
		return nil, err
	}

	return &pb.RemoveWishListResp{
		Success: true,
	}, nil
}
