package user

import (
	"context"
	"heart-trip/app/travel/cmd/rpc/travel"
	"heart-trip/common/ctxdata"

	"heart-trip/app/usercenter/cmd/api/internal/svc"
	"heart-trip/app/usercenter/cmd/api/internal/types"

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

	_, err := l.svcCtx.TravelRpc.AddWishList(l.ctx, &travel.AddWishListReq{
		UserId:     userId,
		HomestayId: req.HomestayId,
	})
	if err != nil {
		return nil, err
	}

	//var resp types.Homestay
	//_ = copier.Copy(&resp, wishListResp.Homestay)

	return &types.AddWishListResp{
		Success: true,
	}, nil
}
