package logic

import (
	"context"
	"time"

	"github.com/zeromicro/go-zero/core/stores/sqlx"

	"heart-trip/app/travel/cmd/rpc/internal/svc"
	"heart-trip/app/travel/cmd/rpc/pb"

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
	userHistory, err := l.svcCtx.UserHistoryModel.FindOneByUserIdAndHistoryId(l.ctx, in.UserId, in.HistoryId)
	history, err := l.svcCtx.HistoryModel.FindOne(l.ctx, in.HistoryId)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(l.ctx, 5*time.Second)
	defer cancel()
	// 放事务里面，失败一个就回滚
	err = l.svcCtx.UserHistoryModel.Transact(ctx, func(ctx context.Context, session sqlx.Session) error {
		err = l.svcCtx.HistoryModel.DeleteSoft(ctx, session, history)
		if err != nil {
			return err
		}
		err = l.svcCtx.UserHistoryModel.Delete(ctx, session, userHistory.Id)
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
