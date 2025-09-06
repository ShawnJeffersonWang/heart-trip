package logic

import (
	"context"
	"errors"
	"fmt"
	"heart-trip/app/travel/model"
	upb "heart-trip/app/usercenter/cmd/rpc/pb"
	"heart-trip/common/globalkey"
	"strconv"
	"strings"

	"github.com/go-redis/redis/v8"
	"github.com/jinzhu/copier"

	"heart-trip/app/travel/cmd/rpc/internal/svc"
	"heart-trip/app/travel/cmd/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type QueryBlogOfFollowLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewQueryBlogOfFollowLogic(ctx context.Context, svcCtx *svc.ServiceContext) *QueryBlogOfFollowLogic {
	return &QueryBlogOfFollowLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// QueryBlogOfFollow 查询关注的博客
func (l *QueryBlogOfFollowLogic) QueryBlogOfFollow(in *pb.QueryBlogOfFollowRequest) (*pb.QueryBlogOfFollowResponse, error) {
	userId := in.UserId
	key := globalkey.FeedKey + strconv.FormatInt(userId, 10)

	// ZREVRANGEBYSCORE key max min LIMIT offset count
	zRangeBy := &redis.ZRangeBy{
		Min:    "0",
		Max:    strconv.FormatInt(in.Max, 10),
		Offset: int64(in.Offset),
		Count:  2, // 假设每次查询2条
	}
	typedTuples, err := l.svcCtx.RedisClient.ZRevRangeByScoreWithScores(l.ctx, key, zRangeBy).Result()
	if err != nil || len(typedTuples) == 0 {
		return &pb.QueryBlogOfFollowResponse{
			Code:    200,
			Data:    &pb.ScrollResult{},
			Message: "成功",
		}, nil
	}

	var ids []int64
	var minTime int64
	var os int
	for _, tuple := range typedTuples {
		id, err := strconv.ParseInt(tuple.Member.(string), 10, 64)
		if err != nil {
			continue
		}
		ids = append(ids, id)
		timeScore := int64(tuple.Score)

		if timeScore == minTime {
			os++
		} else {
			minTime = timeScore
			os = 1
		}
	}

	if len(ids) == 0 {
		return &pb.QueryBlogOfFollowResponse{
			Code:    200,
			Data:    &pb.ScrollResult{},
			Message: "成功",
		}, nil
	}

	idStr := strings.Trim(strings.Replace(fmt.Sprint(ids), " ", ",", -1), "[]")
	var blogs []model.Blog
	if err := l.svcCtx.DB.Where("id IN ?", ids).Order(fmt.Sprintf("FIELD(id, %s)", idStr)).Find(&blogs).Error; err != nil {
		return nil, Fail("查询失败")
	}

	for i := range blogs {
		l.queryBlogUser(&blogs[i])
		l.isBlogLiked(in.UserId, &blogs[i])
	}
	var list []*pb.Blog
	_ = copier.Copy(&list, blogs)
	scrollResult := &pb.ScrollResult{
		List:    list,
		Offset:  int32(os),
		MinTime: minTime,
	}

	return &pb.QueryBlogOfFollowResponse{
		Code:    200,
		Data:    scrollResult,
		Message: "成功",
	}, nil
}

// queryBlogUser 私有方法：查询博客作者信息
func (l *QueryBlogOfFollowLogic) queryBlogUser(blog *model.Blog) {
	//var user model.User
	user, err := l.svcCtx.Usercenter.GetUserInfo(l.ctx, &upb.GetUserInfoReq{Id: blog.UserID})
	if err != nil {
		// 处理错误，例如设置默认值
		blog.Name = "未知用户"
		blog.Icon = ""
		return
	}
	blog.Name = user.User.Nickname
	blog.Icon = user.User.Avatar
}

// isBlogLiked 私有方法：判断当前用户是否点赞
func (l *QueryBlogOfFollowLogic) isBlogLiked(userId int64, blog *model.Blog) {
	key := globalkey.BlogLikedKey + strconv.FormatInt(blog.ID, 10)
	score, err := l.svcCtx.RedisClient.ZScore(l.ctx, key, strconv.FormatInt(userId, 10)).Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		// 处理错误，如记录日志
		l.Logger.Errorf("查询是否点赞失败: %v", err)
		return
	}
	blog.IsLike = score != 0
}
