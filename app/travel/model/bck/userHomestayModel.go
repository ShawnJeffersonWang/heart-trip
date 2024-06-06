package bck

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ UserHomestayModel = (*customUserHomestayModel)(nil)

type (
	// UserHomestayModel is an interface to be customized, add more methods here,
	// and implement the added methods in customUserHomestayModel.
	UserHomestayModel interface {
		userHomestayModel
	}

	customUserHomestayModel struct {
		*defaultUserHomestayModel
	}
)

// NewUserHomestayModel returns a model for the database table.
func NewUserHomestayModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) UserHomestayModel {
	return &customUserHomestayModel{
		defaultUserHomestayModel: newUserHomestayModel(conn, c, opts...),
	}
}
