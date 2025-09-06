package user

import (
	"context"
	"heart-trip/app/usercenter/cmd/rpc/usercenter"
	"heart-trip/common/ctxdata"

	"github.com/jinzhu/copier"

	"heart-trip/app/usercenter/cmd/api/internal/svc"
	"heart-trip/app/usercenter/cmd/api/internal/types"

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
			Nickname: req.Nickname,
			Sex:      req.Sex,
			Avatar:   req.Avatar,
			Info:     req.Info,
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
