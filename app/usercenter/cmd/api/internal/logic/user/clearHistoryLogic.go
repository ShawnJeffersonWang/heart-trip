package user

import (
	"context"
	"homestay/app/travel/cmd/rpc/travel"
	"homestay/common/ctxdata"

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

func (l *ClearHistoryLogic) ClearHistory(req *types.ClearHistoryReq) (*types.ClearHistoryResp, error) {
	// todo: add your logic here and delete this line
	userId := ctxdata.GetUidFromCtx(l.ctx)
	_, err := l.svcCtx.TravelRpc.ClearHistory(l.ctx, &travel.ClearHistoryReq{
		UserId: userId,
	})
	if err != nil {
		return nil, err
	}

	return &types.ClearHistoryResp{
		Success: true,
	}, nil
}
