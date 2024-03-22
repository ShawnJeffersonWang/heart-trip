package user

import (
	"context"
	"github.com/jinzhu/copier"
	"golodge/app/usercenter/cmd/rpc/usercenter"
	"golodge/common/ctxdata"

	"golodge/app/usercenter/cmd/api/internal/svc"
	"golodge/app/usercenter/cmd/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateUserInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateUserInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateUserInfoLogic {
	return &UpdateUserInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateUserInfoLogic) UpdateUserInfo(req *types.UpdateUserInfoReq) (*types.UpdateUserInfoResp, error) {
	// todo: add your logic here and delete this line
	userId := ctxdata.GetUidFromCtx(l.ctx)

	updateUserInfoResp, err := l.svcCtx.UsercenterRpc.UpdateUserInfo(l.ctx, &usercenter.UpdateUserInfoReq{
		User: &usercenter.User{
			Id:       userId,
			Nickname: req.UserInfo.Nickname,
			Sex:      req.UserInfo.Sex,
			Avatar:   req.UserInfo.Avatar,
			Info:     req.UserInfo.Info,
		},
	})
	if err != nil {
		return nil, err
	}

	var userInfo types.User
	_ = copier.Copy(&userInfo, updateUserInfoResp)

	return &types.UpdateUserInfoResp{
		Success: true,
	}, nil
}
