package logic

import (
	"context"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"golodge/app/travel/cmd/rpc/internal/svc"
	"golodge/app/travel/cmd/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type AdminDeleteHomestayLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAdminDeleteHomestayLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminDeleteHomestayLogic {
	return &AdminDeleteHomestayLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *AdminDeleteHomestayLogic) AdminDeleteHomestay(in *pb.AdminDeleteHomestayReq) (*pb.AdminDeleteHomestayResp, error) {
	// todo: add your logic here and delete this line
	homestay, err := l.svcCtx.HomestayModel.FindOne(l.ctx, in.HomestayId)
	homestayActivity, err := l.svcCtx.HomestayActivityModel.FindOneByDataId(l.ctx, in.HomestayId)
	err = l.svcCtx.HomestayModel.Trans(l.ctx, func(ctx context.Context, session sqlx.Session) error {
		err := l.svcCtx.HomestayModel.DeleteSoft(ctx, session, homestay)
		if err != nil {
			return err
		}
		err = l.svcCtx.HomestayActivityModel.DeleteSoft(ctx, session, homestayActivity)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &pb.AdminDeleteHomestayResp{}, nil
}
