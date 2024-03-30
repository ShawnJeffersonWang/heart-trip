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
)

var ErrHomestayAlreadyAdded = xerr.NewErrMsg("homestay has been added in wishList")

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
	homestay, err := l.svcCtx.HomestayModel.FindOne(l.ctx, in.HomestayId)
	if err != nil {
		return nil, err
	}
	// FindOne里面查询不能带缓存会导致再次添加相同的名宿id可以成功，会出现缓存不一致的问题，这可以算在面试中遇到的一个问题
	// 然后查找的时候需要查UserId和HomestayId, 而不能只查HomestayId
	_, err = l.svcCtx.UserHomestayModel.FindOneByUserIdAndHomestayId(l.ctx, in.UserId, in.HomestayId)
	if err == nil {
		return nil, errors.Wrapf(ErrHomestayAlreadyAdded,
			"homestay has been added in wishlist homestayId:%d, err:%v", in.HomestayId, err)
	}
	// bug: 插入无法成功，
	_, err = l.svcCtx.UserHomestayModel.Insert(l.ctx, &model.UserHomestay{
		UserId:     in.UserId,
		HomestayId: in.HomestayId,
	})
	if err != nil {
		return nil, err
	}
	var wishHomestay pb.Homestay
	_ = copier.Copy(&wishHomestay, homestay)
	//Homestay: &pb.Homestay{
	//	Id:          homestay.Id,
	//	Cover:       homestay.Cover,
	//	Title:       homestay.Title,
	//	Intro:       homestay.Intro,
	//	Location:    homestay.Location,
	//	RatingStars: homestay.RatingStars,
	//	PriceBefore: homestay.PriceBefore,
	//	PriceAfter:  homestay.PriceAfter,
	//	RowState:    homestay.RowState,
	//},
	return &pb.AddWishListResp{
		Homestay: &wishHomestay,
	}, nil
}
