package logic

import (
	"context"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"golodge/app/travel/cmd/rpc/internal/svc"
	"golodge/app/travel/cmd/rpc/pb"
	"golodge/app/travel/model"
	"golodge/common/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

var ErrHomestayMismatch = xerr.NewErrMsg("homestay mismatch")

type DeleteHomestayLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteHomestayLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteHomestayLogic {
	return &DeleteHomestayLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeleteHomestayLogic) DeleteHomestay(in *pb.DeleteHomestayReq) (*pb.DeleteHomestayResp, error) {
	// todo: add your logic here and delete this line
	homestay, err := l.svcCtx.HomestayModel.FindOne(l.ctx, in.HomestayId)

	homestayActivity, err := l.svcCtx.HomestayActivityModel.FindOneByDataId(l.ctx, in.HomestayId)
	userHomestay, err := l.svcCtx.UserHomestayModel.FindOneByUserIdAndHomestayId(l.ctx, in.UserId, in.HomestayId)
	history, err := l.svcCtx.HistoryModel.FindOneByHomestayIdAndUserId(l.ctx, in.HomestayId, in.UserId)
	var userHistory *model.UserHistory
	if history != nil {
		userHistory, err = l.svcCtx.UserHistoryModel.FindOne(l.ctx, in.UserId, history.Id)
		if err != nil {
			return nil, err
		}
	}
	// bug 不能写FindOneByHomestayIdAndUserId id本身就是唯一的
	if in.UserId != homestay.HostId {
		return nil, errors.Wrapf(xerr.NewErrMsg(" not authorization "), "userId : %d ", in.UserId)
	}
	err = l.svcCtx.HomestayModel.Trans(l.ctx, func(ctx context.Context, session sqlx.Session) error {
		err := l.svcCtx.HomestayModel.DeleteSoft(ctx, session, homestay)
		if err != nil {
			return err
		}
		err = l.svcCtx.HomestayActivityModel.DeleteSoft(ctx, session, homestayActivity)
		if err != nil {
			return err
		}
		if userHomestay != nil {
			err = l.svcCtx.UserHomestayModel.DeleteSoft(ctx, session, userHomestay)
			if err != nil {
				return err
			}
		}
		if userHistory != nil {
			err = l.svcCtx.UserHistoryModel.DeleteSoft(ctx, session, userHistory)
			if err != nil {
				return err
			}
		}
		if history != nil {
			err = l.svcCtx.HistoryModel.DeleteSoft(ctx, session, history)
			if err != nil {
				return err
			}
			return nil
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &pb.DeleteHomestayResp{}, nil
}
