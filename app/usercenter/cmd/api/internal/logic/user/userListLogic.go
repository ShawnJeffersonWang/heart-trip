package user

import (
	"context"
	"github.com/jinzhu/copier"
	"golodge/app/usercenter/cmd/rpc/pb"

	"golodge/app/usercenter/cmd/api/internal/svc"
	"golodge/app/usercenter/cmd/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserListLogic {
	return &UserListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserListLogic) UserList(req *types.UserListReq) (*types.UserListResp, error) {
	// todo: add your logic here and delete this line
	userListResp, err := l.svcCtx.UsercenterRpc.UserList(l.ctx, &pb.UserListReq{})
	if err != nil {
		return nil, err
	}
	var resp []types.User
	_ = copier.Copy(&resp, userListResp.List)
	return &types.UserListResp{
		List: resp,
	}, nil
}
