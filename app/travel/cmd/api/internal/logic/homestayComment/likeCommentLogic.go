package homestayComment

import (
	"context"
	"golodge/app/travel/cmd/rpc/travel"

	"golodge/app/travel/cmd/api/internal/svc"
	"golodge/app/travel/cmd/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type LikeCommentLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLikeCommentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LikeCommentLogic {
	return &LikeCommentLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LikeCommentLogic) LikeComment(req *types.LikeCommentReq) (*types.LikeCommentResp, error) {
	// todo: add your logic here and delete this line
	_, err := l.svcCtx.TravelRpc.LikeComment(l.ctx, &travel.LikeCommentReq{
		CommentId: req.CommentId,
	})
	if err != nil {
		return nil, err
	}

	// 这样Success字段后不用加,
	return &types.LikeCommentResp{Success: true}, nil
}
