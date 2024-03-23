package homestayComment

import (
	"context"
	"golodge/app/travel/cmd/api/internal/svc"
	"golodge/app/travel/cmd/api/internal/types"
	"golodge/app/travel/cmd/rpc/pb"
	"golodge/app/travel/cmd/rpc/travel"
	"golodge/common/ctxdata"

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
	userId := ctxdata.GetUidFromCtx(l.ctx)
	homestayComment := pb.HomestayComment{
		HomestayId: req.HomestayId,
		Content:    req.Content,
		Star:       req.Star,
		UserId:     userId,
		Nickname:   req.Nickname,
		Avatar:     req.Avatar,
		ImageUrls:  req.ImageUrls,
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
