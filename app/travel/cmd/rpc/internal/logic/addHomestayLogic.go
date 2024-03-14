package logic

import (
	"context"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"homestay/app/travel/cmd/rpc/internal/svc"
	"homestay/app/travel/cmd/rpc/pb"
	"homestay/app/travel/cmd/rpc/travel"
	"homestay/app/travel/model"

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

	if err := l.svcCtx.HomestayModel.Trans(l.ctx, func(ctx context.Context, session sqlx.Session) error {
		homestay := model.Homestay{
			Title:       in.Homestay.Title,
			Cover:       in.Homestay.Cover,
			Intro:       in.Homestay.Intro,
			Location:    in.Homestay.Location,
			UserId:      in.Homestay.UserId,
			PriceBefore: in.Homestay.PriceBefore,
		}
		l.svcCtx.HomestayModel.Insert(ctx, session, &homestay)
		return nil
	}); err != nil {
		return nil, err
	}

	return &travel.AddHomestayResp{}, nil
}
