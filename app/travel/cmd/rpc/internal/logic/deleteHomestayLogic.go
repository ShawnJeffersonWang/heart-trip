package logic

import (
	"context"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"golodge/app/travel/cmd/rpc/internal/svc"
	"golodge/app/travel/cmd/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

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
	_, err := l.svcCtx.HomestayModel.FindOneByUserId(l.ctx, in.UserId)
	//if in.UserId != homestay.UserId {
	//	return nil, errors.Wrapf(xerr.NewErrCode(xerr.DB_ERROR), " No Authorization err , userId : %d ", in.UserId)
	//}
	err = l.svcCtx.HomestayModel.Trans(l.ctx, func(context context.Context, session sqlx.Session) error {
		err := l.svcCtx.HomestayModel.Delete(l.ctx, session, in.HomestayId)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &pb.DeleteHomestayResp{}, nil
}
