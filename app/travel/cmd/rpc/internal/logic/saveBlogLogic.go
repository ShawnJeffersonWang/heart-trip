package logic

import (
	"context"
	"fmt"
	"heart-trip/app/travel/model"
	upb "heart-trip/app/usercenter/cmd/rpc/pb"
	"heart-trip/common/globalkey"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"

	"heart-trip/app/travel/cmd/rpc/internal/svc"
	"heart-trip/app/travel/cmd/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type SaveBlogLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSaveBlogLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SaveBlogLogic {
	return &SaveBlogLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// SaveBlog 保存博客
func (l *SaveBlogLogic) SaveBlog(in *pb.SaveBlogRequest) (*pb.SaveBlogResponse, error) {
	blog := model.Blog{
		UserID: in.UserId,
		// 其他字段赋值...
	}

	if err := l.svcCtx.DB.Create(&blog).Error; err != nil {
		return nil, Fail("新增笔记失败!")
	}

	// 查询笔记作者的所有粉丝
	followUserIDResponse, err := l.svcCtx.Usercenter.QueryFollowsByFollowUserID(l.ctx, &upb.QueryFollowsByFollowUserIDRequest{FollowUserId: in.UserId})
	if err != nil {
		return nil, Fail("获取粉丝失败")
	}
	follows := followUserIDResponse.Data

	// 推送笔记id给所有粉丝
	for _, follow := range follows {
		key := globalkey.FeedKey + strconv.FormatInt(follow.UserId, 10)
		if err := l.svcCtx.RedisClient.ZAdd(l.ctx, key, &redis.Z{
			Score:  float64(time.Now().Unix()),
			Member: strconv.FormatInt(blog.ID, 10),
		}).Err(); err != nil {
			// 可以选择忽略或记录日志
			l.Logger.Errorf("推送博客ID到Redis失败: %v", err)
			continue
		}
	}

	return &pb.SaveBlogResponse{
		Code:    200,
		Data:    blog.ID,
		Message: "新增成功",
	}, nil
}

// getUser 私有方法：获取当前用户
//func (l *BlogServiceLogic) getUser() *UserDTO {
//	// 实现从上下文中获取用户信息的方法
//	// 例如，可以从JWT或其他认证方式中提取用户信息
//	// 这里假设通过上下文存储用户信息
//	val := l.ctx.Value("user")
//	if val == nil {
//		return nil
//	}
//	user, ok := val.(*UserDTO)
//	if !ok {
//		return nil
//	}
//	return user
//}

// Fail 辅助函数：返回失败结果
func Fail(msg string) error {
	return fmt.Errorf(msg)
}
