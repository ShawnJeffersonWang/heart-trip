package logic

import (
	"context"
	"github.com/zeromicro/go-zero/core/stores/sqlx"

	"golodge/app/travel/cmd/rpc/internal/svc"
	"golodge/app/travel/cmd/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type RemoveHistoryLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRemoveHistoryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RemoveHistoryLogic {
	return &RemoveHistoryLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *RemoveHistoryLogic) RemoveHistory(in *pb.RemoveHistoryReq) (*pb.RemoveHistoryResp, error) {
	// todo: add your logic here and delete this line
	userHistory, err := l.svcCtx.UserHistoryModel.FindOne(l.ctx, in.UserId, in.HistoryId)
	history, err := l.svcCtx.HistoryModel.FindOne(l.ctx, in.HistoryId)
	// 放事务里面，失败一个就回滚
	err = l.svcCtx.UserHistoryModel.Trans(l.ctx, func(ctx context.Context, session sqlx.Session) error {
		err = l.svcCtx.HistoryModel.DeleteSoft(ctx, session, history)
		if err != nil {
			return err
		}
		err = l.svcCtx.UserHistoryModel.DeleteSoft(ctx, session, userHistory)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &pb.RemoveHistoryResp{}, nil
}
