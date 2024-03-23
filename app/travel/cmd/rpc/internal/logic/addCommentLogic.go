package logic

import (
	"context"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"golodge/app/travel/model"

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
		HomestayId: in.HomestayComment.HomestayId,
		Content:    in.HomestayComment.Content,
		Star:       in.HomestayComment.Star,
		UserId:     in.HomestayComment.UserId,
	}
	if err := l.svcCtx.HomestayCommentModel.Trans(l.ctx, func(ctx context.Context, session sqlx.Session) error {
		l.svcCtx.HomestayCommentModel.Insert(l.ctx, session, &homestayComment)
		return nil
	}); err != nil {
		return nil, err
	}
	return &pb.AddCommentResp{}, nil
}
