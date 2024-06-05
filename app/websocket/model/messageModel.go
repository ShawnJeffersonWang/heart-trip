package model

import (
	"context"
	"database/sql"
	"errors"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"time"
)

type Message struct {
	Id         int64     `db:"id"`
	FromUserId int64     `db:"from_user_id"`
	ToUserId   int64     `db:"to_user_id"`
	Content    string    `db:"content"`
	DelState   int64     `db:"del_state"`
	CreateTime time.Time `db:"create_time"`
}

type MessageModel interface {
	FindByUserId(ctx context.Context, userId int64) ([]Message, error)
	Insert(ctx context.Context, data *Message) (sql.Result, error)
}

type defaultMessageModel struct {
	conn  sqlx.SqlConn
	table string
}

func NewMessageModel(conn sqlx.SqlConn) MessageModel {
	return &defaultMessageModel{
		conn:  conn,
		table: "message",
	}
}

func (m *defaultMessageModel) Insert(ctx context.Context, data *Message) (sql.Result, error) {
	query := `INSERT INTO ` + m.table + ` (from_user_id, to_user_id, content) 
              VALUES (?, ?, ?)`
	return m.conn.ExecCtx(ctx, query, data.FromUserId, data.ToUserId, data.Content)
}

func (m *defaultMessageModel) FindByUserId(ctx context.Context, userId int64) ([]Message, error) {
	query := `SELECT id, from_user_id, to_user_id, content, del_state, create_time 
              FROM ` + m.table + ` 
              WHERE to_user_id = ? AND del_state = 0 
              ORDER BY create_time DESC`
	var messages []Message
	err := m.conn.QueryRowsCtx(ctx, &messages, query, userId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return messages, nil
}
