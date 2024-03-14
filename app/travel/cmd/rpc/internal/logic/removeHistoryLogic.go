package logic

import (
	"context"

	"homestay/app/travel/cmd/rpc/internal/svc"
	"homestay/app/travel/cmd/rpc/pb"

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

	return &pb.RemoveHistoryResp{}, nil
}
