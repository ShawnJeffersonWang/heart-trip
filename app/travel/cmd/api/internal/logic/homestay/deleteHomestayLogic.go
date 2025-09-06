package homestay

import (
	"context"
	"heart-trip/app/travel/cmd/api/internal/svc"
	"heart-trip/app/travel/cmd/api/internal/types"
	"heart-trip/app/travel/cmd/rpc/pb"
	"heart-trip/common/ctxdata"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteHomestayLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteHomestayLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteHomestayLogic {
	return &DeleteHomestayLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteHomestayLogic) DeleteHomestay(req *types.DeleteHomestayReq) (*types.DeleteHomestayResp, error) {
	// todo: add your logic here and delete this line
	userId := ctxdata.GetUidFromCtx(l.ctx)
	// 写一个轮子直接调用, 方便后期扩展
	//searchHistoryResp, err := l.svcCtx.TravelRpc.SearchHistory(l.ctx, &pb.SearchHistoryReq{
	//	UserId:     userId,
	//	HomestayId: req.HomestayId,
	//})
	//if searchHistoryResp != nil {
	//	_, err = l.svcCtx.TravelRpc.RemoveHistory(l.ctx, &pb.RemoveHistoryReq{
	//		UserId:    userId,
	//		HistoryId: searchHistoryResp.History.Id,
	//	})
	//	if err != nil {
	//		return nil, err
	//	}
	//}
	//// RPC写好了直接调用, 越写越快，良性循环, 就像搭积木一样
	//_, err = l.svcCtx.TravelRpc.RemoveWishList(l.ctx, &pb.RemoveWishListReq{
	//	UserId:     userId,
	//	HomestayId: req.HomestayId,
	//})
	//if err != nil {
	//	return nil, err
	//}
	_, err := l.svcCtx.TravelRpc.DeleteHomestay(l.ctx, &pb.DeleteHomestayReq{
		UserId:     userId,
		HomestayId: req.HomestayId,
	})
	if err != nil {
		return nil, err
	}

	return &types.DeleteHomestayResp{
		Success: true,
	}, nil
}
