package user

import (
	"context"
	"github.com/jinzhu/copier"
	"homestay/app/travel/cmd/rpc/travel"
	"homestay/common/ctxdata"

	"homestay/app/usercenter/cmd/api/internal/svc"
	"homestay/app/usercenter/cmd/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddWishListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAddWishListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddWishListLogic {
	return &AddWishListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AddWishListLogic) AddWishList(req *types.AddWishListReq) (*types.AddWishListResp, error) {
	// todo: add your logic here and delete this line
	userId := ctxdata.GetUidFromCtx(l.ctx)

	wishListResp, err := l.svcCtx.TravelRpc.AddWishList(l.ctx, &travel.AddWishListReq{
		UserId:     userId,
		HomestayId: req.Id,
	})
	if err != nil {
		return nil, err
	}

	var resp types.Homestay
	_ = copier.Copy(&resp, wishListResp.Homestay)

	return &types.AddWishListResp{
		Homestay: resp,
	}, nil
}
