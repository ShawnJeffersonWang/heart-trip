package genModel

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ UserHistoryModel = (*customUserHistoryModel)(nil)

type (
	// UserHistoryModel is an interface to be customized, add more methods here,
	// and implement the added methods in customUserHistoryModel.
	UserHistoryModel interface {
		userHistoryModel
	}

	customUserHistoryModel struct {
		*defaultUserHistoryModel
	}
)

// NewUserHistoryModel returns a model for the database table.
func NewUserHistoryModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) UserHistoryModel {
	return &customUserHistoryModel{
		defaultUserHistoryModel: newUserHistoryModel(conn, c, opts...),
	}
}
