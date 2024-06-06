package user

import (
	"context"
	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"
	"golodge/app/travel/cmd/rpc/travel"
	"golodge/app/usercenter/cmd/api/internal/svc"
	"golodge/app/usercenter/cmd/api/internal/types"
	"golodge/common/ctxdata"
)

type WishListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewWishListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *WishListLogic {
	return &WishListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *WishListLogic) WishList(req *types.WishListReq) (*types.WishListResp, error) {
	// todo: add your logic here and delete this line
	// bug: 之所以会出现服务器问题是因为svcCtx没有对TravelRpc初始化，不建议在svcCtx加别的服务的model,建议使用rpc的方式调用别的服务
	// 这样就很神奇了，让我感受到了微服务的魅力
	userId := ctxdata.GetUidFromCtx(l.ctx)

	wishListResp, err := l.svcCtx.TravelRpc.WishList(l.ctx, &travel.WishListReq{
		UserId: userId,
	})
	if err != nil {
		return nil, err
	}

	var resp []types.Homestay
	_ = copier.Copy(&resp, wishListResp.List)

	return &types.WishListResp{
		List: resp,
	}, nil
}
