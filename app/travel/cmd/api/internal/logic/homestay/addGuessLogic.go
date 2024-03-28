package homestay

import (
	"context"
	"golodge/app/travel/cmd/rpc/pb"

	"golodge/app/travel/cmd/api/internal/svc"
	"golodge/app/travel/cmd/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddGuessLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAddGuessLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddGuessLogic {
	return &AddGuessLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AddGuessLogic) AddGuess(req *types.AddGuessReq) (*types.AddGuessResp, error) {
	// todo: add your logic here and delete this line
	_, err := l.svcCtx.TravelRpc.AddGuess(l.ctx, &pb.AddGuessReq{
		HomestayId: req.HomestayId,
	})
	if err != nil {
		return nil, err
	}
	return &types.AddGuessResp{
		Success: true,
	}, nil
}
