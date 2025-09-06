package logic

import (
	"context"
	"errors"
	"golodge/app/travel/model"
	upb "golodge/app/usercenter/cmd/rpc/pb"
	"golodge/common/globalkey"
	"strconv"

	"github.com/go-redis/redis/v8"
	"github.com/jinzhu/copier"

	"golodge/app/travel/cmd/rpc/internal/svc"
	"golodge/app/travel/cmd/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type QueryBlogByIdLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewQueryBlogByIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *QueryBlogByIdLogic {
	return &QueryBlogByIdLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// QueryBlogById 根据ID查询博客
func (l *QueryBlogByIdLogic) QueryBlogById(in *pb.QueryBlogByIdRequest) (*pb.QueryBlogByIdResponse, error) {
	var blog model.Blog
	if err := l.svcCtx.DB.First(&blog, in.Id).Error; err != nil {
		return nil, Fail("笔记不存在！")
	}

	l.queryBlogUser(&blog)
	l.isBlogLiked(in.UserId, &blog)
	var data *pb.Blog
	_ = copier.Copy(&data, &blog)
	return &pb.QueryBlogByIdResponse{
		Code:    200,
		Data:    data,
		Message: "成功",
	}, nil
}

// queryBlogUser 私有方法：查询博客作者信息
func (l *QueryBlogByIdLogic) queryBlogUser(blog *model.Blog) {
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
func (l *QueryBlogByIdLogic) isBlogLiked(userId int64, blog *model.Blog) {
	key := globalkey.BlogLikedKey + strconv.FormatInt(blog.ID, 10)
	score, err := l.svcCtx.RedisClient.ZScore(l.ctx, key, strconv.FormatInt(userId, 10)).Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		// 处理错误，如记录日志
		l.Logger.Errorf("查询是否点赞失败: %v", err)
		return
	}
	blog.IsLike = score != 0
}
