package user

import (
	"context"
	"golodge/app/travel/cmd/rpc/travel"
	"golodge/common/ctxdata"

	"golodge/app/usercenter/cmd/api/internal/svc"
	"golodge/app/usercenter/cmd/api/internal/types"

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

func (l *RemoveHistoryLogic) RemoveHistory(req *types.RemoveHistoryReq) (*types.RemoveHistoryResp, error) {
	// todo: add your logic here and delete this line
	userId := ctxdata.GetUidFromCtx(l.ctx)
	_, err := l.svcCtx.TravelRpc.RemoveHistory(l.ctx, &travel.RemoveHistoryReq{
		UserId:    userId,
		HistoryId: req.HistoryId,
	})
	if err != nil {
		return nil, err
	}

	return &types.RemoveHistoryResp{
		Success: true,
	}, nil
}
