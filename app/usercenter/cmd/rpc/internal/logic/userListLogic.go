package logic

import (
	"context"
	"github.com/Masterminds/squirrel"
	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/mr"
	"golodge/app/usercenter/model"

	"golodge/app/usercenter/cmd/rpc/internal/svc"
	"golodge/app/usercenter/cmd/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUserListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserListLogic {
	return &UserListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UserListLogic) UserList(in *pb.UserListReq) (*pb.UserListResp, error) {
	// todo: add your logic here and delete this line
	whereBuilder := l.svcCtx.UserModel.SelectBuilder().Where(squirrel.Eq{})
	users, err := l.svcCtx.UserModel.FindAll(l.ctx, whereBuilder, "id")
	if err != nil {
		return nil, err
	}
	var resp []*pb.User
	if len(users) > 0 {
		mr.MapReduceVoid(func(source chan<- any) {
			for _, user := range users {
				source <- user.Id
			}
		}, func(item any, writer mr.Writer[*model.User], cancel func(error)) {
			id := item.(int64)
			user, err := l.svcCtx.UserModel.FindOne(l.ctx, id)
			if err != nil {
				logx.WithContext(l.ctx).Errorf("获取活动数据失败: id %d, err %v", id, err)
			}
			writer.Write(user)
		}, func(pipe <-chan *model.User, cancel func(error)) {
			for user := range pipe {
				var tyUser pb.User
				_ = copier.Copy(&tyUser, user)
				resp = append(resp, &tyUser)
			}
		})
	}

	return &pb.UserListResp{
		List: resp,
	}, nil
}
