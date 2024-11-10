package logic

import (
	"context"
	"fmt"
	"golodge/app/usercenter/model"
	"golodge/common/globalkey"
	"strconv"
	"strings"

	"golodge/app/travel/cmd/rpc/internal/svc"
	"golodge/app/travel/cmd/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type QueryBlogLikesLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewQueryBlogLikesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *QueryBlogLikesLogic {
	return &QueryBlogLikesLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// QueryBlogLikes 查询某篇博客的点赞用户
func (l *QueryBlogLikesLogic) QueryBlogLikes(in *pb.QueryBlogLikesRequest) (*pb.QueryBlogLikesResponse, error) {
	key := globalkey.BlogLikedKey + strconv.FormatInt(in.Id, 10)
	top5, err := l.svcCtx.RedisClient.ZRange(l.ctx, key, 0, 4).Result()
	if err != nil || len(top5) == 0 {
		return &pb.QueryBlogLikesResponse{
			Code:    200,
			Data:    []*pb.UserDTO{},
			Message: "成功",
		}, nil
	}

	var ids []int64
	for _, uidStr := range top5 {
		uid, err := strconv.ParseInt(uidStr, 10, 64)
		if err == nil {
			ids = append(ids, uid)
		}
	}

	if len(ids) == 0 {
		return &pb.QueryBlogLikesResponse{
			Code:    200,
			Data:    []*pb.UserDTO{},
			Message: "成功",
		}, nil
	}

	var users []model.User
	idStr := strings.Trim(strings.Replace(fmt.Sprint(ids), " ", ",", -1), "[]")
	if err := l.svcCtx.DB.Where("id IN ?", ids).Order(fmt.Sprintf("FIELD(id, %s)", idStr)).Find(&users).Error; err != nil {
		return nil, Fail("查询失败")
	}

	var userDTOs []*pb.UserDTO
	for _, user := range users {
		userDTOs = append(userDTOs, &pb.UserDTO{
			Id:       user.Id,
			NickName: user.Nickname,
			Icon:     user.Avatar,
		})
	}

	return &pb.QueryBlogLikesResponse{
		Code:    200,
		Data:    userDTOs,
		Message: "成功",
	}, nil
}
