package logic

import (
	"context"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"golodge/app/usercenter/cmd/rpc/usercenter"
	"golodge/app/usercenter/model"
	"golodge/common/xerr"

	"golodge/app/usercenter/cmd/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateUserInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateUserInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateUserInfoLogic {
	return &UpdateUserInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateUserInfoLogic) UpdateUserInfo(in *usercenter.UpdateUserInfoReq) (*usercenter.UpdateUserInfoResp, error) {
	// todo: add your logic here and delete this line
	user, err := l.svcCtx.UserModel.FindOne(l.ctx, in.User.Id)
	if err != nil && err != model.ErrNotFound {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DB_ERROR), "GetUserInfo find user db err , id:%d , err:%v", in.User.Id, err)
	}
	if user == nil {
		return nil, errors.Wrapf(ErrUserNoExistsError, "id:%d", in.User.Id)
	}

	user.Nickname = in.User.Nickname
	user.Sex = in.User.Id
	user.Avatar = in.User.Avatar
	user.Info = in.User.Info
	if err := l.svcCtx.UserModel.Trans(l.ctx, func(ctx context.Context, session sqlx.Session) error {
		l.svcCtx.UserModel.Update(ctx, session, user)
		return nil
	}); err != nil {
		return nil, err
	}

	return &usercenter.UpdateUserInfoResp{}, nil
}
