package homestayComment

import (
	"context"
	"heart-trip/app/travel/cmd/api/internal/svc"
	"heart-trip/app/travel/cmd/api/internal/types"
	"heart-trip/app/travel/cmd/rpc/pb"
	"heart-trip/app/travel/cmd/rpc/travel"
	"heart-trip/common/ctxdata"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddCommentLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAddCommentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddCommentLogic {
	return &AddCommentLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AddCommentLogic) AddComment(req *types.AddCommentReq) (*types.AddCommentResp, error) {
	// todo: add your logic here and delete this line
	// 为什么感觉到微服务越写越快，因为轮子造好了，rpc调用就完事了
	userId := ctxdata.GetUidFromCtx(l.ctx)
	homestayComment := pb.HomestayComment{
		HomestayId:     req.HomestayId,
		CommentTime:    req.CommentTime,
		Content:        req.Content,
		Star:           req.Star,
		TidyRating:     req.TidyRating,
		TrafficRating:  req.TrafficRating,
		SecurityRating: req.SecurityRating,
		FoodRating:     req.FoodRating,
		CostRating:     req.CostRating,
		UserId:         userId,
		Nickname:       req.Nickname,
		Avatar:         req.Avatar,
		ImageUrls:      req.ImageUrls,
	}
	_, err := l.svcCtx.TravelRpc.AddComment(l.ctx, &travel.AddCommentReq{
		HomestayComment: &homestayComment,
	})
	if err != nil {
		return nil, err
	}

	return &types.AddCommentResp{
		Success: true,
	}, nil
}
