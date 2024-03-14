package genModel

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ HistoryHomestayModel = (*customHistoryHomestayModel)(nil)

type (
	// HistoryHomestayModel is an interface to be customized, add more methods here,
	// and implement the added methods in customHistoryHomestayModel.
	HistoryHomestayModel interface {
		historyHomestayModel
	}

	customHistoryHomestayModel struct {
		*defaultHistoryHomestayModel
	}
)

// NewHistoryHomestayModel returns a model for the database table.
func NewHistoryHomestayModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) HistoryHomestayModel {
	return &customHistoryHomestayModel{
		defaultHistoryHomestayModel: newHistoryHomestayModel(conn, c, opts...),
	}
}
