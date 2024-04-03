package logic

import (
	"context"
	"github.com/jinzhu/copier"

	"golodge/app/travel/cmd/rpc/internal/svc"
	"golodge/app/travel/cmd/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type SearchHistoryLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSearchHistoryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SearchHistoryLogic {
	return &SearchHistoryLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SearchHistoryLogic) SearchHistory(in *pb.SearchHistoryReq) (*pb.SearchHistoryResp, error) {
	// todo: add your logic here and delete this line
	history, err := l.svcCtx.HistoryModel.FindOneByHomestayIdAndUserId(l.ctx, in.HomestayId, in.UserId)
	if err != nil {
		return nil, err
	}
	var pbHistory pb.History
	_ = copier.Copy(&pbHistory, history)
	return &pb.SearchHistoryResp{
		History: &pbHistory,
	}, nil
}
