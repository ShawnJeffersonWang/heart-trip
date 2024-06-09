package model

import (
	"context"
	"fmt"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"time"
)

type UserHistory struct {
	Id         int64     `db:"id"`
	HistoryId  int64     `db:"history_id"`
	UserId     int64     `db:"user_id"`
	DelState   int64     `db:"del_state"`
	Version    int64     `db:"version"`
	DeleteTime time.Time `db:"delete_time"`
}

type UserHistoryModel interface {
	Insert(ctx context.Context, history *UserHistory) (int64, error)
	Delete(ctx context.Context, session sqlx.Session, id int64) error
	FindOneByUserIdAndHistoryId(ctx context.Context, userId int64, historyId int64) (*UserHistory, error)
	FindUserHistories(ctx context.Context, userId int64, pageNum, pageSize int64) ([]UserHistory, error)
	FindByUserId(ctx context.Context, userId int64) ([]*UserHistory, error)
	Transact(ctx context.Context, fn func(ctx context.Context, session sqlx.Session) error) error
}
type defaultUserHistoryModel struct {
	conn  sqlx.SqlConn
	table string
}

func NewUserHistoryModel(conn sqlx.SqlConn, cache cache.CacheConf) UserHistoryModel {
	return &defaultUserHistoryModel{
		conn:  conn,
		table: "`user_history`",
	}
}

// 插入用户历史记录
func (m *defaultUserHistoryModel) Insert(ctx context.Context, history *UserHistory) (int64, error) {
	query := "INSERT INTO user_history (history_id, user_id) VALUES (?, ?)"
	result, err := m.conn.ExecCtx(ctx, query, history.HistoryId, history.UserId)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

// 逻辑删除用户历史记录
func (m *defaultUserHistoryModel) Delete(ctx context.Context, session sqlx.Session, id int64) error {
	query := "UPDATE user_history SET del_state = 1, delete_time = ? WHERE id = ?"
	_, err := m.conn.ExecCtx(ctx, query, time.Now(), id)
	return err
}

func (m *defaultUserHistoryModel) FindOneByUserIdAndHistoryId(ctx context.Context, userId int64, historyId int64) (*UserHistory, error) {
	query := `SELECT id, history_id, user_id, del_state, version, delete_time 
              FROM ` + m.table + ` 
              WHERE user_id = ? AND history_id = ? AND del_state = 0 LIMIT 1`
	var userHistory UserHistory
	err := m.conn.QueryRowCtx(ctx, &userHistory, query, userId, historyId)
	switch err {
	case nil:
		return &userHistory, nil
	case sqlc.ErrNotFound:
		return nil, nil
	default:
		return nil, err
	}
}

func (m *defaultUserHistoryModel) FindUserHistories(ctx context.Context, userId int64, pageNum, pageSize int64) ([]UserHistory, error) {
	offset := (pageNum - 1) * pageSize
	query := fmt.Sprintf("SELECT * FROM user_history WHERE user_id = ? AND del_state = 0 ORDER BY id DESC LIMIT ? OFFSET ?")
	var histories []UserHistory
	err := m.conn.QueryRowsCtx(ctx, &histories, query, userId, pageSize, offset)
	if err != nil {
		return nil, err
	}
	return histories, nil
}

func (m *defaultUserHistoryModel) FindByUserId(ctx context.Context, userId int64) ([]*UserHistory, error) {
	query := `SELECT id, history_id, user_id, del_state, version, delete_time FROM ` + m.table + ` WHERE user_id = ?`
	var histories []*UserHistory
	err := m.conn.QueryRowsCtx(ctx, &histories, query, userId)
	if err != nil {
		return nil, err
	}
	return histories, nil
}

func (m *defaultUserHistoryModel) Transact(ctx context.Context, fn func(ctx context.Context, session sqlx.Session) error) error {
	return m.conn.TransactCtx(ctx, fn)
}
