package logic

import (
	"context"
	"golodge/app/travel/model"
	"golodge/app/usercenter/cmd/rpc/internal/svc"
	"golodge/app/usercenter/cmd/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type QueryFollowsByFollowUserIDLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewQueryFollowsByFollowUserIDLogic(ctx context.Context, svcCtx *svc.ServiceContext) *QueryFollowsByFollowUserIDLogic {
	return &QueryFollowsByFollowUserIDLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *QueryFollowsByFollowUserIDLogic) QueryFollowsByFollowUserID(in *pb.QueryFollowsByFollowUserIDRequest) (*pb.QueryFollowsByFollowUserIDResponse, error) {
	var follows []model.Follow
	if err := l.svcCtx.DB.WithContext(l.ctx).Where("follow_user_id = ?", in.FollowUserId).Find(follows).Error; err != nil {
		// 返回错误响应
		return &pb.QueryFollowsByFollowUserIDResponse{
			Code:    500,
			Message: "获取粉丝失败",
		}, nil
	}

	// 转换 model.Follow 到 Protobuf Follow
	var pbFollows []*pb.Follow
	for _, f := range follows {
		pbFollows = append(pbFollows, &pb.Follow{
			UserId:       f.UserID,
			FollowUserId: f.FollowUserID,
		})
	}

	// 返回成功响应
	return &pb.QueryFollowsByFollowUserIDResponse{
		Code:    200,
		Data:    pbFollows,
		Message: "成功",
	}, nil
}
