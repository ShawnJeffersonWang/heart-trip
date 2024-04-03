package homestay

import (
	"context"
	"golodge/app/travel/cmd/api/internal/svc"
	"golodge/app/travel/cmd/api/internal/types"
	"golodge/app/travel/cmd/rpc/travel"
	"golodge/common/ctxdata"
	"golodge/common/xerr"

	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
)

type HomestayDetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewHomestayDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) HomestayDetailLogic {
	return HomestayDetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *HomestayDetailLogic) HomestayDetail(req types.HomestayDetailReq) (*types.HomestayDetailResp, error) {

	userId := ctxdata.GetUidFromCtx(l.ctx)
	homestayResp, err := l.svcCtx.TravelRpc.HomestayDetail(l.ctx, &travel.HomestayDetailReq{
		HomestayId: req.HomestayId,
		UserId:     userId,
	})
	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrMsg("get homestay detail fail"),
			" get homestay detail db err , id : %d , err : %v ", req.HomestayId, err)
	}

	//var typeHomestay types.Homestay
	//if homestayResp != nil {
	//	_ = copier.Copy(&typeHomestay, homestayResp)
	//}

	return &types.HomestayDetailResp{
		Id:                 homestayResp.Id,
		Title:              homestayResp.Title,
		RatingStars:        homestayResp.RatingStars,
		CommentCount:       homestayResp.CommentCount,
		TitleTags:          homestayResp.TitleTags,
		BannerUrls:         homestayResp.BannerUrls,
		Latitude:           homestayResp.Latitude,
		Longitude:          homestayResp.Longitude,
		Location:           homestayResp.Location,
		Facilities:         homestayResp.Facilities,
		Area:               homestayResp.Area,
		RoomConfig:         homestayResp.RoomConfig,
		CleanVideo:         homestayResp.CleanVideo,
		HostAvatar:         homestayResp.HostAvatar,
		HostNickname:       homestayResp.HostNickname,
		PriceBefore:        homestayResp.PriceBefore,
		PriceAfter:         homestayResp.PriceAfter,
		HostId:             homestayResp.HostId,
		HomestayBusinessId: homestayResp.HomestayBusinessId,
	}, nil
}
