package logic

import (
	"context"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"golodge/app/travel/cmd/rpc/internal/svc"
	"golodge/app/travel/cmd/rpc/pb"
	"golodge/app/travel/cmd/rpc/travel"
	"golodge/app/travel/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddHomestayLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAddHomestayLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddHomestayLogic {
	return &AddHomestayLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *AddHomestayLogic) AddHomestay(in *pb.AddHomestayReq) (*pb.AddHomestayResp, error) {
	// todo: add your logic here and delete this line

	//_, err := l.svcCtx.HomestayModel.FindOne(l.ctx, in.Homestay.Id)
	//if err == nil {
	//	return nil, errors.Wrapf(ErrHomestayAlreadyAdded,
	//		"homestay has been added in homestayList homestayId:%d,err:%v", in.Homestay.Id, err)
	//}

	if err := l.svcCtx.HomestayModel.Trans(l.ctx, func(ctx context.Context, session sqlx.Session) error {
		homestay := model.Homestay{
			Title:       in.Homestay.Title,
			Cover:       in.Homestay.Cover,
			Intro:       in.Homestay.Intro,
			Location:    in.Homestay.Location,
			UserId:      in.Homestay.UserId,
			PriceBefore: in.Homestay.PriceBefore,
		}

		_, err := l.svcCtx.HomestayModel.Insert(ctx, session, &homestay)
		if err != nil {
			return err
		}

		homestayActivity := model.HomestayActivity{
			DelState:  0,
			RowType:   "preferredHomestay",
			DataId:    homestay.Id,
			RowStatus: 1,
			Version:   0,
		}

		_, err = l.svcCtx.HomestayActivityModel.Insert(ctx, session, &homestayActivity)
		if err != nil {
			return err
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return &travel.AddHomestayResp{}, nil
}
