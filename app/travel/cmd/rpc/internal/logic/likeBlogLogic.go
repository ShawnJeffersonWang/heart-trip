package logic

import (
	"context"
	"errors"
	"github.com/go-redis/redis/v8"
	"golodge/app/travel/model"
	"golodge/common/globalkey"
	"gorm.io/gorm"
	"strconv"
	"time"

	"golodge/app/travel/cmd/rpc/internal/svc"
	"golodge/app/travel/cmd/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type LikeBlogLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewLikeBlogLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LikeBlogLogic {
	return &LikeBlogLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// LikeBlog 点赞或取消点赞
func (l *LikeBlogLogic) LikeBlog(in *pb.LikeBlogRequest) (*pb.LikeBlogResponse, error) {
	userId := in.UserId
	key := globalkey.BlogLikedKey + strconv.FormatInt(in.Id, 10)

	score, err := l.svcCtx.RedisClient.ZScore(l.ctx, key, strconv.FormatInt(userId, 10)).Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		return nil, Fail("点赞失败")
	}

	if score == 0 {
		// 未点赞，执行点赞操作
		if err := l.svcCtx.DB.Model(&model.Blog{}).Where("id = ?", in.Id).UpdateColumn("liked", gorm.Expr("liked + ?", 1)).Error; err != nil {
			return nil, Fail("点赞失败")
		}
		// 保存用户到Redis的有序集合
		if err := l.svcCtx.RedisClient.ZAdd(l.ctx, key, &redis.Z{
			Score:  float64(time.Now().Unix()),
			Member: strconv.FormatInt(userId, 10),
		}).Err(); err != nil {
			return nil, Fail("点赞失败")
		}
	} else {
		// 已点赞，执行取消点赞操作
		if err := l.svcCtx.DB.Model(&model.Blog{}).Where("id = ?", in.Id).UpdateColumn("liked", gorm.Expr("liked - ?", 1)).Error; err != nil {
			return nil, Fail("取消点赞失败")
		}
		// 从Redis的有序集合中移除用户
		if err := l.svcCtx.RedisClient.ZRem(l.ctx, key, strconv.FormatInt(userId, 10)).Err(); err != nil {
			return nil, Fail("取消点赞失败")
		}
	}

	return &pb.LikeBlogResponse{
		Code:    200,
		Message: "操作成功",
	}, nil
}
