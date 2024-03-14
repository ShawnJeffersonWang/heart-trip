package logic

import (
	"context"

	"homestay/app/travel/cmd/rpc/internal/svc"
	"homestay/app/travel/cmd/rpc/pb"

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
	// todo: add your logic here and delete this line

	return &pb.ClearHistoryResp{}, nil
}
