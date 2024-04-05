package logic

import (
	"context"
	"github.com/jinzhu/copier"
	"github.com/pkg/errors"
	"golodge/app/travel/cmd/rpc/internal/svc"
	"golodge/app/travel/cmd/rpc/pb"
	"golodge/app/travel/model"
	"golodge/common/xerr"

	"github.com/zeromicro/go-zero/core/logx"
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

func (l *HomestayDetailWithoutLoginLogic) HomestayDetailWithoutLogin(in *pb.HomestayDetailReq) (*pb.HomestayDetailResp, error) {
	homestay, err := l.svcCtx.HomestayModel.FindOne(l.ctx, in.HomestayId)
	if err != nil && err != model.ErrNotFound {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DB_ERROR), " HomestayDetail db err , id : %d ", in.HomestayId)
	}
	var pbHomestay pb.Homestay
	if homestay != nil {
		_ = copier.Copy(&pbHomestay, homestay)
	}
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
