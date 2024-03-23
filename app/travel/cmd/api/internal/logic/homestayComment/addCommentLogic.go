package homestayComment

import (
	"context"
	"golodge/app/travel/cmd/api/internal/svc"
	"golodge/app/travel/cmd/api/internal/types"
	"golodge/app/travel/cmd/rpc/pb"
	"golodge/app/travel/cmd/rpc/travel"

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
	homestayComment := pb.HomestayComment{
		HomestayId: req.HomestayComment.HomestayId,
		Content:    req.HomestayComment.Content,
		Star:       req.HomestayComment.Star,
		UserId:     req.HomestayComment.UserId,
		Nickname:   req.HomestayComment.Nickname,
		Avatar:     req.HomestayComment.Avatar,
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
