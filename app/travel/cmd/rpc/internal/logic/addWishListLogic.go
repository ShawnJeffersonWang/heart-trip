package logic

import (
	"context"
	"github.com/zeromicro/go-zero/core/logx"
	"homestay/app/travel/cmd/rpc/internal/svc"
	"homestay/app/travel/cmd/rpc/pb"
	"homestay/app/travel/model"
)

type AddWishListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAddWishListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddWishListLogic {
	return &AddWishListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *AddWishListLogic) AddWishList(in *pb.AddWishListReq) (*pb.AddWishListResp, error) {
	// todo: add your logic here and delete this line
	homestay, _ := l.svcCtx.HomestayModel.FindOne(l.ctx, in.HomestayId)

	_, err := l.svcCtx.UserHomestayModel.Insert(l.ctx, &model.UserHomestay{
		UserId:     in.UserId,
		HomestayId: in.HomestayId,
	})
	if err != nil {
		return nil, err
	}

	return &pb.AddWishListResp{
		Homestay: &pb.Homestay{
			Id:          homestay.Id,
			Cover:       homestay.Cover,
			Title:       homestay.Title,
			Intro:       homestay.Intro,
			Location:    homestay.Location,
			RatingStars: homestay.RatingStars,
			PriceBefore: homestay.PriceBefore,
			PriceAfter:  homestay.PriceAfter,
			RowState:    homestay.RowState,
		},
	}, nil
}
