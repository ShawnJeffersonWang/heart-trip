package logic

import (
	"context"
	"heart-trip/app/travel/model"
	"time"

	"heart-trip/app/travel/cmd/rpc/internal/svc"
	"heart-trip/app/travel/cmd/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddGuessLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAddGuessLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddGuessLogic {
	return &AddGuessLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *AddGuessLogic) AddGuess(in *pb.AddGuessReq) (*pb.AddGuessResp, error) {
	// todo: add your logic here and delete this line
	homestay, err := l.svcCtx.HomestayModel.FindOne(l.ctx, in.HomestayId)
	if err != nil {
		return nil, err
	}
	guess := model.Guess{
		Id:          0,
		HomestayId:  homestay.Id,
		PriceAfter:  homestay.PriceAfter,
		PriceBefore: homestay.PriceBefore,
		Cover:       homestay.Cover,
		Location:    homestay.Location,
		Title:       homestay.Title,
		IsCollected: 0,
		UdateTime:   time.Now(),
		CreateTime:  time.Now(),
	}
	_, err = l.svcCtx.GuessModel.Insert(l.ctx, &guess)
	if err != nil {
		return nil, err
	}

	return &pb.AddGuessResp{}, nil
}
