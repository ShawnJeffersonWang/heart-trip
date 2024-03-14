package user

import (
	"context"

	"homestay/app/usercenter/cmd/api/internal/svc"
	"homestay/app/usercenter/cmd/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ClearHistoryLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewClearHistoryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ClearHistoryLogic {
	return &ClearHistoryLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ClearHistoryLogic) ClearHistory(req *types.ClearHistoryReq) (resp *types.ClearHistoryResp, err error) {
	// todo: add your logic here and delete this line

	return
}
