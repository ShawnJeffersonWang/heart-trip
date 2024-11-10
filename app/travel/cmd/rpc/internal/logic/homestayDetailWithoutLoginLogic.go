package logic

import (
	"context"
	"fmt"
	"github.com/zeromicro/go-zero/core/logx"
	"golodge/app/travel/cmd/rpc/internal/svc"
	"golodge/app/travel/cmd/rpc/pb"
	"golodge/app/travel/model"
	"golodge/common/globalkey"
)

type HomestayDetailWithoutLoginLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewHomestayDetailWithoutLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *HomestayDetailWithoutLoginLogic {
	return &HomestayDetailWithoutLoginLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// getById 从数据库中根据ID获取Shop
func (l *HomestayDetailWithoutLoginLogic) getById(id any) (any, error) {
	// 类型断言为 int64（根据具体需求调整）
	shopID, ok := id.(int64)
	if !ok {
		return nil, fmt.Errorf("invalid id type")
	}

	ctx := context.Background() // 或者传入上下文
	shop, err := l.svcCtx.ShopModel.FindOne(ctx, shopID)
	if err != nil {
		return nil, err
	}
	return shop, nil
}

func (l *HomestayDetailWithoutLoginLogic) HomestayDetailWithoutLogin(in *pb.HomestayDetailReq) (*pb.HomestayDetailResp, error) {
	data, err := l.svcCtx.CacheClient.QueryWithPassThrough(
		l.ctx,
		globalkey.CacheShopKey,
		in.HomestayId,
		l.getById,
		globalkey.CacheShopTtl,
	)
	if err != nil {
		return nil, Fail("店铺不存在！")
	}

	// 类型断言为 *model.Shop
	homestay, ok := data.(*model.Homestay)
	if !ok || homestay == nil {
		return nil, Fail("店铺不存在！")
	}
	//homestay, err := l.svcCtx.HomestayModel.FindOne(l.ctx, in.HomestayId)
	//if err != nil && err != model.ErrNotFound {
	//	return nil, errors.Wrapf(xerr.NewErrCode(xerr.DB_ERROR), " HomestayDetail db err , id : %d ", in.HomestayId)
	//}
	//var pbHomestay pb.Homestay
	//if homestay != nil {
	//	_ = copier.Copy(&pbHomestay, homestay)
	//}
	return &pb.HomestayDetailResp{
		Id:                 homestay.Id,
		Title:              homestay.Title,
		RatingStars:        homestay.RatingStars,
		CommentCount:       homestay.CommentCount,
		TitleTags:          homestay.TitleTags,
		BannerUrls:         homestay.BannerUrls,
		Latitude:           homestay.Latitude,
		Longitude:          homestay.Longitude,
		Location:           homestay.Location,
		Facilities:         homestay.Facilities,
		Area:               homestay.Area,
		RoomConfig:         homestay.RoomConfig,
		CleanVideo:         homestay.CleanVideo,
		HostAvatar:         homestay.HostAvatar,
		HostNickname:       homestay.HostNickname,
		PriceBefore:        homestay.PriceBefore,
		PriceAfter:         homestay.PriceAfter,
		HostId:             homestay.HostId,
		HomestayBusinessId: homestay.HomestayBusinessId,
	}, nil
}
