package homestay

import (
	"context"

	"homestay/app/travel/cmd/api/internal/svc"
	"homestay/app/travel/cmd/api/internal/types"
	"homestay/common/xerr"

	"github.com/jinzhu/copier"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
)

type GuessListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGuessListLogic(ctx context.Context, svcCtx *svc.ServiceContext) GuessListLogic {
	return GuessListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GuessListLogic) GuessList(req types.GuessListReq) (*types.GuessListResp, error) {
	var resp []types.Homestay

	list, err := l.svcCtx.HomestayModel.FindPageListByIdDESC(l.ctx, l.svcCtx.HomestayModel.SelectBuilder(), 0, 5)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DB_ERROR), "GuessList db err req : %+v , err : %v", req, err)
	}

	if len(list) > 0 {
		for _, homestay := range list {
			var typeHomestay types.Homestay
			_ = copier.Copy(&typeHomestay, homestay)

			resp = append(resp, typeHomestay)
		}
	}

	return &types.GuessListResp{
		List: resp,
	}, nil
}
