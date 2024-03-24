package model

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ GuessModel = (*customGuessModel)(nil)

type (
	// GuessModel is an interface to be customized, add more methods here,
	// and implement the added methods in customGuessModel.
	GuessModel interface {
		guessModel
	}

	customGuessModel struct {
		*defaultGuessModel
	}
)

// NewGuessModel returns a model for the database table.
func NewGuessModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) GuessModel {
	return &customGuessModel{
		defaultGuessModel: newGuessModel(conn, c, opts...),
	}
}
