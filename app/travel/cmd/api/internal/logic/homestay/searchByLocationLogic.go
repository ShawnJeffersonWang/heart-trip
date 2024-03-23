package homestay

import (
	"context"
	"github.com/jinzhu/copier"
	"golodge/app/travel/cmd/rpc/pb"

	"golodge/app/travel/cmd/api/internal/svc"
	"golodge/app/travel/cmd/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type SearchByLocationLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSearchByLocationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SearchByLocationLogic {
	return &SearchByLocationLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SearchByLocationLogic) SearchByLocation(req *types.SearchByLocationReq) (*types.SearchByLocationResp, error) {
	// todo: add your logic here and delete this line
	searchByLocationResp, err := l.svcCtx.TravelRpc.SearchByLocation(l.ctx, &pb.SearchByLocationReq{
		Location: req.Location,
	})
	if err != nil {
		return nil, err
	}

	var resp []types.Homestay
	if searchByLocationResp != nil {
		_ = copier.Copy(&resp, searchByLocationResp.List)
	}
	return &types.SearchByLocationResp{
		List: resp,
	}, nil
}
