package logic

import (
	"context"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"golodge/app/travel/model"
	"strconv"

	"golodge/app/travel/cmd/rpc/internal/svc"
	"golodge/app/travel/cmd/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddCommentLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAddCommentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddCommentLogic {
	return &AddCommentLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *AddCommentLogic) AddComment(in *pb.AddCommentReq) (*pb.AddCommentResp, error) {
	// todo: add your logic here and delete this line
	homestayComment := model.HomestayComment{
		HomestayId:     in.HomestayComment.HomestayId,
		CommentTime:    in.HomestayComment.CommentTime,
		Content:        in.HomestayComment.Content,
		Star:           in.HomestayComment.Star,
		UserId:         in.HomestayComment.UserId,
		Nickname:       in.HomestayComment.Nickname,
		Avatar:         in.HomestayComment.Avatar,
		ImageUrls:      in.HomestayComment.ImageUrls,
		CostRating:     in.HomestayComment.CostRating,
		TrafficRating:  in.HomestayComment.TrafficRating,
		TidyRating:     in.HomestayComment.TidyRating,
		SecurityRating: in.HomestayComment.SecurityRating,
		FoodRating:     in.HomestayComment.FoodRating,
	}
	homestay, err := l.svcCtx.HomestayModel.FindOne(l.ctx, in.HomestayComment.HomestayId)
	if err != nil {
		return nil, err
	}
	homestay.CommentCount++
	ratingStar, _ := strconv.ParseFloat(in.HomestayComment.Star, 64)
	cnt := float64(homestay.CommentCount)
	homestay.RatingStars = (homestay.RatingStars + ratingStar) / cnt
	if err := l.svcCtx.HomestayCommentModel.Trans(l.ctx, func(ctx context.Context, session sqlx.Session) error {
		_, err := l.svcCtx.HomestayCommentModel.Insert(ctx, session, &homestayComment)
		if err != nil {
			return err
		}
		// 分别调用了HomestayComment和Homestay的服务，将他们组装成一个事务，要么成功要么回滚
		_, err = l.svcCtx.HomestayModel.Update(ctx, session, homestay)
		if err != nil {
			return err
		}
		return nil
	}); err != nil {
		return nil, err
	}
	return &pb.AddCommentResp{}, nil
}
