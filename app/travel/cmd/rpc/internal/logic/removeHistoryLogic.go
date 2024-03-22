package logic

import (
	"context"

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

	err := l.svcCtx.HistoryHomestayModel.Delete(l.ctx, in.UserId, in.HistoryId)
	if err != nil {
		return nil, err
	}

	err = l.svcCtx.HistoryModel.Delete(l.ctx, in.HistoryId)
	if err != nil {
		return nil, err
	}
	return &pb.RemoveHistoryResp{}, nil
}
