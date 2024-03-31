package homestay

import (
	"context"
	"golodge/app/travel/cmd/api/internal/svc"
	"golodge/app/travel/cmd/api/internal/types"
	"golodge/app/travel/cmd/rpc/travel"
	"golodge/app/usercenter/cmd/rpc/pb"
	"golodge/common/ctxdata"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddHomestayLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAddHomestayLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddHomestayLogic {
	return &AddHomestayLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AddHomestayLogic) AddHomestay(req *types.AddHomestayReq) (*types.AddHomestayResp, error) {
	// todo: add your logic here and delete this line
	userId := ctxdata.GetUidFromCtx(l.ctx)
	getUserInfoResp, err := l.svcCtx.UsercenterRpc.GetUserInfo(l.ctx, &pb.GetUserInfoReq{
		Id: userId,
	})
	if err != nil {
		return nil, err
	}
	_, err = l.svcCtx.TravelRpc.AddHomestay(l.ctx, &travel.AddHomestayReq{
		Title:        req.Title,
		TitleTags:    req.TitleTags,
		BannerUrls:   req.BannerUrls,
		Latitude:     req.Latitude,
		Longitude:    req.Longitude,
		Facilities:   req.Facilities,
		Area:         req.Area,
		RoomConfig:   req.RoomConfig,
		CleanVideo:   req.CleanVideo,
		PriceBefore:  req.PriceBefore,
		PriceAfter:   req.PriceAfter,
		HostId:       userId,
		HostAvatar:   getUserInfoResp.User.Avatar,
		HostNickname: getUserInfoResp.User.Nickname,
	})
	if err != nil {
		return nil, err
	}

	return &types.AddHomestayResp{
		Success: true,
	}, nil
}
