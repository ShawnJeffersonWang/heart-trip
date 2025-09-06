package logic

import (
	"context"
	"heart-trip/common/globalkey"
	"strconv"

	"heart-trip/app/travel/cmd/rpc/internal/svc"
	"heart-trip/app/travel/cmd/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateShopLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateShopLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateShopLogic {
	return &UpdateShopLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateShopLogic) UpdateShop(in *pb.UpdateShopRequest) (*pb.UpdateShopResponse, error) {
	shop := in.Homestay
	if shop.Id == 0 {
		return &pb.UpdateShopResponse{
			Code:    400,
			Message: "店铺id不能为空",
		}, nil
	}

	// 1. 更新数据库
	err := l.svcCtx.DB.Save(&shop).Error
	if err != nil {
		l.Errorf("更新店铺失败: %v", err)
		return &pb.UpdateShopResponse{
			Code:    500,
			Message: "更新店铺失败",
		}, nil
	}

	// 2. 删除缓存
	cacheKey := globalkey.CacheShopKey + strconv.FormatInt(shop.Id, 10)
	err = l.svcCtx.RedisClient.Del(l.ctx, cacheKey).Err()
	if err != nil {
		l.Errorf("删除缓存失败: %v", err)
		// 根据需求决定是否返回错误，通常可以记录日志后继续
	}

	return &pb.UpdateShopResponse{
		Code:    200,
		Message: "更新成功",
	}, nil
}
