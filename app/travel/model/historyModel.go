package model

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ HistoryModel = (*customHistoryModel)(nil)

type (
	// HistoryModel is an interface to be customized, add more methods here,
	// and implement the added methods in customHistoryModel.
	HistoryModel interface {
		historyModel
	}

	customHistoryModel struct {
		*defaultHistoryModel
	}
)

// NewHistoryModel returns a model for the database table.
func NewHistoryModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) HistoryModel {
	return &customHistoryModel{
		defaultHistoryModel: newHistoryModel(conn, c, opts...),
	}
}
