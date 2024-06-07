package homestay

import (
	"context"
	"golodge/app/travel/cmd/api/internal/svc"
	"golodge/app/travel/cmd/api/internal/types"
	"golodge/common/xerr"

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
	var resp []types.Guess

	// bug: 改了前面的没改后面的
	list, err := l.svcCtx.GuessModel.FindPageListByIdDESC(l.ctx, l.svcCtx.GuessModel.SelectBuilder(), 0, 5)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DB_ERROR), "GuessList db err req : %+v , err : %v", req, err)
	}

	if len(list) > 0 {
		for _, guess := range list {
			//l.svcCtx.UserHomestayModel.CheckIfExists(l.ctx,)
			var typeHomestay types.Guess
			_ = copier.Copy(&typeHomestay, guess)

			resp = append(resp, typeHomestay)
		}
	}

	return &types.GuessListResp{
		List: resp,
	}, nil
}
