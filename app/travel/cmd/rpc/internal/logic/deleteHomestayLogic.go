package logic

import (
	"context"
	"github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"golodge/app/travel/cmd/rpc/internal/svc"
	"golodge/app/travel/cmd/rpc/pb"
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
	// bug 不能写FindOneByHomestayIdAndUserId id本身就是唯一的
	if in.UserId != homestay.HostId {
		return nil, errors.Wrapf(xerr.NewErrMsg(" not authorization "), "userId : %d ", in.UserId)
	}
	homestayActivity, err := l.svcCtx.HomestayActivityModel.FindOneByDataId(l.ctx, in.HomestayId)
	whereBuilder := l.svcCtx.UserHomestayModel.SelectBuilder().Where(squirrel.Eq{
		"homestay_id": in.HomestayId,
	})
	userHomestays, err := l.svcCtx.UserHomestayModel.FindAll(l.ctx, whereBuilder, "id desc")
	// bug: 不能直接用userId查，结果不对
	//userHomestay, err := l.svcCtx.UserHomestayModel.FindOneByUserIdAndHomestayId(l.ctx, in.UserId, in.HomestayId)
	// bug: 不能直接用userId查
	//history, err := l.svcCtx.HistoryModel.FindOneByHomestayIdAndUserId(l.ctx, in.HomestayId, in.UserId)
	whereBuilder = l.svcCtx.HistoryModel.SelectBuilder().Where(squirrel.Eq{
		"homestay_id": in.HomestayId,
	})
	histories, err := l.svcCtx.HistoryModel.FindAll(l.ctx, whereBuilder, "id desc")
	err = l.svcCtx.HomestayModel.Trans(l.ctx, func(ctx context.Context, session sqlx.Session) error {
		err := l.svcCtx.HomestayModel.DeleteSoft(ctx, session, homestay)
		if err != nil {
			return err
		}
		err = l.svcCtx.HomestayActivityModel.DeleteSoft(ctx, session, homestayActivity)
		if err != nil {
			return err
		}
		for _, userHomestay := range userHomestays {
			err := l.svcCtx.UserHomestayModel.DeleteSoft(ctx, session, userHomestay)
			if err != nil {
				return err
			}
		}
		for _, history := range histories {
			err := l.svcCtx.HistoryModel.DeleteSoft(ctx, session, history)
			if err != nil {
				return err
			}

			userHistory, err := l.svcCtx.UserHistoryModel.FindOne(ctx, in.UserId, history.Id)
			if err != nil {
				return err
			}
			err = l.svcCtx.UserHistoryModel.DeleteSoft(ctx, session, userHistory)
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &pb.DeleteHomestayResp{}, nil
}
