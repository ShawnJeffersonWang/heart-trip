package user

import (
	"context"

	"homestay/app/usercenter/cmd/api/internal/svc"
	"homestay/app/usercenter/cmd/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type RemoveHistoryLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRemoveHistoryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RemoveHistoryLogic {
	return &RemoveHistoryLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RemoveHistoryLogic) RemoveHistory(req *types.RemoveHistoryReq) (resp *types.RemoveHistoryResp, err error) {
	// todo: add your logic here and delete this line

	return
}
