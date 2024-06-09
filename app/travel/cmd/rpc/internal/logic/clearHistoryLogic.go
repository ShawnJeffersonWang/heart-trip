package logic

import (
	"context"
	"github.com/Masterminds/squirrel"
	"github.com/zeromicro/go-zero/core/stores/sqlx"

	"golodge/app/travel/cmd/rpc/internal/svc"
	"golodge/app/travel/cmd/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type ClearHistoryLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewClearHistoryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ClearHistoryLogic {
	return &ClearHistoryLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ClearHistoryLogic) ClearHistory(in *pb.ClearHistoryReq) (*pb.ClearHistoryResp, error) {
	//whereBuilder := l.svcCtx.UserHistoryModel.SelectBuilder().Where(squirrel.Eq{
	//	"user_id": in.UserId,
	//})
	userHistories, err := l.svcCtx.UserHistoryModel.FindByUserId(l.ctx, in.UserId)

	whereBuilder := l.svcCtx.HistoryModel.SelectBuilder().Where(squirrel.Eq{
		"user_id": in.UserId,
	})
	histories, err := l.svcCtx.HistoryModel.FindAll(l.ctx, whereBuilder, "")
	err = l.svcCtx.UserHistoryModel.Transact(l.ctx, func(ctx context.Context, session sqlx.Session) error {
		for _, userHistory := range userHistories {
			err := l.svcCtx.UserHistoryModel.Delete(ctx, session, userHistory.Id)
			if err != nil {
				return err
			}
		}
		for _, history := range histories {
			err = l.svcCtx.HistoryModel.DeleteSoft(ctx, session, history)
			if err != nil {
				return err
			}
		}
		return nil
	})
	//err := l.svcCtx.HistoryModel.DeleteAll(l.ctx, in.UserId)
	//err = l.svcCtx.UserHistoryModel.DeleteAll(l.ctx, in.UserId)
	if err != nil {
		return nil, err
	}
	return &pb.ClearHistoryResp{}, nil
}
