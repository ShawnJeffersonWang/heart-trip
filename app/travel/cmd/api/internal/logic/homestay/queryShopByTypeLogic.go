package homestay

import (
	"context"
	"github.com/jinzhu/copier"
	"golodge/app/travel/cmd/rpc/travel"

	"golodge/app/travel/cmd/api/internal/svc"
	"golodge/app/travel/cmd/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type QueryShopByTypeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewQueryShopByTypeLogic query shop by type
func NewQueryShopByTypeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *QueryShopByTypeLogic {
	return &QueryShopByTypeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *QueryShopByTypeLogic) QueryShopByType(req *types.QueryShopByTypeRequest) (*types.QueryShopByTypeResponse, error) {
	queryShopByTypeResponse, err := l.svcCtx.TravelRpc.QueryShopByType(l.ctx, &travel.QueryShopByTypeRequest{
		TypeId:  req.TypeId,
		Current: req.Current,
		X:       req.X,
		Y:       req.Y,
	})
	if err != nil {
		return nil, err
	}
	var data []types.Homestay
	_ = copier.Copy(&data, queryShopByTypeResponse.Data)
	return &types.QueryShopByTypeResponse{
		Code: int(queryShopByTypeResponse.Code),
		Msg:  queryShopByTypeResponse.Msg,
		Data: data,
	}, nil
}
