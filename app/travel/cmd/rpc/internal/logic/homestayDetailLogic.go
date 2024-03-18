package logic

import (
	"context"
	"fmt"
	"github.com/jinzhu/copier"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
	"homestay/app/travel/cmd/rpc/internal/svc"
	"homestay/app/travel/cmd/rpc/pb"
	"homestay/app/travel/model"
	"homestay/common/xerr"
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

	homestay, err := l.svcCtx.HomestayModel.FindOne(l.ctx, in.Id)
	if err != nil && err != model.ErrNotFound {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DB_ERROR), " HomestayDetail db err , id : %d ", in.Id)
	}

	history, err := l.svcCtx.HistoryModel.FindOneByHomestayId(l.ctx, in.Id)
	fmt.Println("history: ", history)
	if history == nil {
		// bug: 这里用的是homestay.UserId, 导致一直是0, 应该用in.UserId
		history = &model.History{
			Title:              homestay.Title,
			HomestayBusinessId: homestay.HomestayBusinessId,
			Intro:              homestay.Intro,
			Cover:              homestay.Cover,
			Location:           homestay.Location,
			PriceAfter:         homestay.PriceAfter,
			PriceBefore:        homestay.PriceBefore,
			RatingStars:        homestay.RatingStars,
			UserId:             in.UserId,
			HomestayId:         homestay.Id,
			CreateTime:         time.Now(),
			UpdateTime:         time.Now(),
		}
		_, err = l.svcCtx.HistoryModel.Insert(l.ctx, history)

		historyTemp, err := l.svcCtx.HistoryModel.FindOneByUserId(l.ctx, in.UserId)
		fmt.Println("history: ", historyTemp)

		historyHomestay := &model.HistoryHomestay{
			HistoryId: historyTemp.Id,
			UserId:    in.UserId,
		}

		_, err = l.svcCtx.HistoryHomestayModel.Insert(l.ctx, historyHomestay)
		if err != nil {
			return nil, err
		}
	}
	var pbHomestay pb.Homestay
	if homestay != nil {
		_ = copier.Copy(&pbHomestay, homestay)
	}

	return &pb.HomestayDetailResp{
		Homestay: &pbHomestay,
	}, nil

}
