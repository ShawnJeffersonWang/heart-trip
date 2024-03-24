package user

import (
	"context"
	"github.com/jinzhu/copier"
	"golodge/app/travel/cmd/rpc/travel"
	"golodge/common/ctxdata"

	"golodge/app/usercenter/cmd/api/internal/svc"
	"golodge/app/usercenter/cmd/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type HistoryListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewHistoryListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *HistoryListLogic {
	return &HistoryListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *HistoryListLogic) HistoryList(req *types.HistoryListReq) (*types.HistoryListResp, error) {
	// todo: add your logic here and delete this line
	userId := ctxdata.GetUidFromCtx(l.ctx)
	historyListResp, err := l.svcCtx.TravelRpc.HistoryList(l.ctx, &travel.HistoryListReq{
		UserId: userId,
	})
	if err != nil {
		return nil, err
	}

	var resp []types.History
	_ = copier.Copy(&resp, historyListResp.HistoryList)
	return &types.HistoryListResp{
		HistoryList: resp,
	}, nil
}
