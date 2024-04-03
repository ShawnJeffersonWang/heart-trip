package logic

import (
	"context"
	"github.com/jinzhu/copier"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
	"golodge/app/travel/cmd/rpc/internal/svc"
	"golodge/app/travel/cmd/rpc/pb"
	"golodge/app/travel/model"
	"golodge/common/xerr"
	"time"
)

type HomestayDetailLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewHomestayDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *HomestayDetailLogic {
	return &HomestayDetailLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// HomestayDetail homestay detail .
func (l *HomestayDetailLogic) HomestayDetail(in *pb.HomestayDetailReq) (*pb.HomestayDetailResp, error) {

	homestay, err := l.svcCtx.HomestayModel.FindOne(l.ctx, in.HomestayId)
	if err != nil && err != model.ErrNotFound {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DB_ERROR), " HomestayDetail db err , id : %d ", in.HomestayId)
	}
	// 缓存不一致, 查缓存发现是以前的历史记录id, 而删不掉真实的, 还是会添加历史记录
	historyTemp, err := l.svcCtx.HistoryModel.FindOneByHomestayIdAndUserId(l.ctx, in.HomestayId, in.UserId)
	// 类似LRU, 找到了先删除, 没找到才新建一个历史记录和建立关联
	if historyTemp != nil {
		_ = l.svcCtx.HistoryModel.Delete(l.ctx, historyTemp.Id)
		_ = l.svcCtx.UserHistoryModel.Delete(l.ctx, in.UserId, historyTemp.Id)
	}
	// bug: 这里用的是homestay.UserId, 导致一直是0, 应该用in.UserId
	history := model.History{
		Title:              homestay.Title,
		HomestayBusinessId: homestay.HomestayBusinessId,
		Intro:              homestay.RoomConfig,
		Cover:              homestay.Cover,
		Location:           homestay.Location,
		PriceAfter:         homestay.PriceAfter,
		PriceBefore:        homestay.PriceBefore,
		RatingStars:        homestay.RatingStars,
		UserId:             in.UserId,
		HomestayId:         homestay.Id,
		CreateTime:         time.Now(),
		LastBrowsingTime:   time.Now(),
	}
	res, err := l.svcCtx.HistoryModel.Insert(l.ctx, &history)

	// bug: 当时没想到可以用&history, 传出参数, 就想的再去通过UserId查历史记录, 其实根本不一定是刚插入的这条历史记录
	//historyTemp, err := l.svcCtx.HistoryModel.FindOneByUserId(l.ctx, in.UserId)
	//fmt.Println("history: ", historyTemp)
	// &history无法获取数据库自增的id, 要使用LastInsertId()
	historyId, _ := res.LastInsertId()
	userHistory := model.UserHistory{
		HistoryId: historyId,
		UserId:    in.UserId,
	}

	_, err = l.svcCtx.UserHistoryModel.Insert(l.ctx, &userHistory)
	if err != nil {
		return nil, err
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
