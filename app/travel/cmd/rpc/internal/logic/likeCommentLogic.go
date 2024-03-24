package logic

import (
	"context"
	"fmt"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"golodge/app/travel/cmd/rpc/internal/svc"
	"golodge/app/travel/cmd/rpc/pb"
)

type LikeCommentLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewLikeCommentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LikeCommentLogic {
	return &LikeCommentLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *LikeCommentLogic) LikeComment(in *pb.LikeCommentReq) (*pb.LikeCommentResp, error) {
	// todo: add your logic here and delete this line
	comment, err := l.svcCtx.HomestayCommentModel.FindOne(l.ctx, in.CommentId)
	if err != nil {
		return nil, err
	}
	comment.LikeCount++
	// bug 不知道为什么不生效
	//homestayComment := model.HomestayComment{
	//	DelState:       0,
	//	HomestayId:     comment.HomestayId,
	//	UserId:         comment.UserId,
	//	Content:        comment.Content,
	//	Star:           comment.Star,
	//	Version:        comment.Version,
	//	Nickname:       comment.Nickname,
	//	Avatar:         comment.Avatar,
	//	ImageUrls:      comment.ImageUrls,
	//	LikeCount:      comment.LikeCount + 1,
	//	CommentTime:    comment.CommentTime,
	//	TidyRating:     comment.TidyRating,
	//	TrafficRating:  comment.TrafficRating,
	//	SecurityRating: comment.SecurityRating,
	//	FoodRating:     comment.FoodRating,
	//	CostRating:     comment.CostRating,
	//}
	if err := l.svcCtx.HomestayCommentModel.Trans(l.ctx, func(context context.Context, session sqlx.Session) error {
		_, err := l.svcCtx.HomestayCommentModel.Update(l.ctx, session, comment)
		if err != nil {
			fmt.Println("更新失败")
			return err
		}
		return nil
	}); err != nil {
		return nil, err
	}
	return &pb.LikeCommentResp{}, nil
}
