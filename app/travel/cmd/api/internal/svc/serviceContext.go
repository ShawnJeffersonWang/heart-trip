package svc

import (
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/zrpc"
	"golodge/app/travel/cmd/api/internal/config"
	"golodge/app/travel/cmd/rpc/travel"
	"golodge/app/travel/model"
	"golodge/app/usercenter/cmd/rpc/usercenter"
)

type ServiceContext struct {
	//local
	Config config.Config

	//rpc
	UsercenterRpc usercenter.Usercenter
	TravelRpc     travel.Travel

	//model
	HomestayModel         model.HomestayModel
	GuessModel            model.GuessModel
	UserHomestayModel     model.UserHomestayModel
	HomestayActivityModel model.HomestayActivityModel
	HomestayBusinessModel model.HomestayBusinessModel
	HomestayCommentModel  model.HomestayCommentModel
}

func NewServiceContext(c config.Config) *ServiceContext {

	sqlConn := sqlx.NewMysql(c.DB.DataSource)

	return &ServiceContext{
		Config: c,

		UsercenterRpc: usercenter.NewUsercenter(zrpc.MustNewClient(c.UsercenterRpcConf)),
		TravelRpc:     travel.NewTravel(zrpc.MustNewClient(c.TravelRpcConf)),

		HomestayModel:         model.NewHomestayModel(sqlConn, c.Cache),
		GuessModel:            model.NewGuessModel(sqlConn, c.Cache),
		UserHomestayModel:     model.NewUserHomestayModel(sqlConn, c.Cache),
		HomestayActivityModel: model.NewHomestayActivityModel(sqlConn, c.Cache),
		HomestayBusinessModel: model.NewHomestayBusinessModel(sqlConn, c.Cache),
		HomestayCommentModel:  model.NewHomestayCommentModel(sqlConn, c.Cache),
	}
}
