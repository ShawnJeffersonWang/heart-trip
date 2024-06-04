package homestay

import (
	"context"
	"github.com/pkg/errors"
	"golodge/app/travel/cmd/rpc/travel"
	"golodge/app/usercenter/cmd/rpc/usercenter"
	"golodge/common/ctxdata"
	"golodge/common/xerr"

	"golodge/app/travel/cmd/api/internal/svc"
	"golodge/app/travel/cmd/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateHomestayLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// update homestay
func NewUpdateHomestayLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateHomestayLogic {
	return &UpdateHomestayLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateHomestayLogic) UpdateHomestay(req *types.UpdateHomestayReq) (*types.UpdateHomestayResp, error) {
	// todo: add your logic here and delete this line
	userId := ctxdata.GetUidFromCtx(l.ctx)
	getUserInfoResp, err := l.svcCtx.UsercenterRpc.GetUserInfo(l.ctx, &usercenter.GetUserInfoReq{
		Id: userId,
	})
	if err != nil {
		return nil, err
	}
	if _, err := l.svcCtx.TravelRpc.UpdateHomestay(l.ctx, &travel.UpdateHomestayReq{
		HomestayId:   req.HomestayId,
		Title:        req.Title,
		TitleTags:    req.TitleTags,
		BannerUrls:   req.BannerUrls,
		Latitude:     req.Latitude,
		Longitude:    req.Longitude,
		Location:     req.Location,
		Facilities:   req.Facilities,
		Area:         req.Area,
		RoomConfig:   req.RoomConfig,
		CleanVideo:   req.CleanVideo,
		PriceBefore:  req.PriceBefore,
		PriceAfter:   req.PriceAfter,
		HostId:       userId,
		HostAvatar:   getUserInfoResp.User.Avatar,
		HostNickname: getUserInfoResp.User.Nickname,
		RowState:     req.RowState,
	}); err != nil {
		return nil, errors.Wrapf(xerr.NewErrMsg(" not authorization "), "userId : %d ", userId)
	}
	return &types.UpdateHomestayResp{Success: true}, nil
}
