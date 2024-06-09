package model

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"time"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type UserHomestay struct {
	Id         int64     `db:"id"`
	UserId     int64     `db:"user_id"`
	HomestayId int64     `db:"homestay_id"`
	DelState   int64     `db:"del_state"`
	Version    int64     `db:"version"`
	DeleteTime time.Time `db:"delete_time"`
}

type UserHomestayModel interface {
	Insert(ctx context.Context, data *UserHomestay) (sql.Result, error)
	FindOne(id int64) (*UserHomestay, error)
	FindAllByHomestayId(ctx context.Context, homestayId int64) ([]UserHomestay, error)
	Update(data *UserHomestay) error
	Delete(id int64) error
	CheckIfExists(ctx context.Context, userId, homestayId int64) (bool, error)
	Favorite(userId, homestayId int64) error
	Unfavorite(userId, homestayId int64) error
	GetFavorites(userId int64, page int64, pageSize int64) ([]*UserHomestay, error)
	UpdateDelState(ctx context.Context, userId, homestayId int64, delState int) error
}

type defaultUserHomestayModel struct {
	conn  sqlx.SqlConn
	table string
}

func NewUserHomestayModel(conn sqlx.SqlConn, cache cache.CacheConf) UserHomestayModel {
	return &defaultUserHomestayModel{
		conn:  conn,
		table: "user_homestay",
	}
}

func (m *defaultUserHomestayModel) Insert(ctx context.Context, data *UserHomestay) (sql.Result, error) {
	query := fmt.Sprintf("INSERT INTO %s (user_id, homestay_id) VALUES (?, ?)", m.table)
	return m.conn.ExecCtx(ctx, query, data.UserId, data.HomestayId)
}

func (m *defaultUserHomestayModel) FindOne(id int64) (*UserHomestay, error) {
	query := `SELECT id, user_id, homestay_id, del_state, version, delete_time FROM ` + m.table + ` WHERE id = ?`
	var data UserHomestay
	err := m.conn.QueryRow(&data, query, id)
	if err != nil {
		return nil, err
	}
	return &data, nil
}

func (m *defaultUserHomestayModel) FindAllByHomestayId(ctx context.Context, homestayId int64) ([]UserHomestay, error) {
	var userHomestays []UserHomestay
	query := `SELECT * FROM ` + m.table + ` WHERE homestay_id = ? AND del_state = 0` // 过滤掉逻辑删除的记录
	err := m.conn.QueryRowsCtx(ctx, &userHomestays, query, homestayId)
	return userHomestays, err
}

func (m *defaultUserHomestayModel) Update(data *UserHomestay) error {
	query := `UPDATE ` + m.table + ` SET user_id = ?, homestay_id = ?, del_state = ?, version = ?, delete_time = ? WHERE id = ?`
	_, err := m.conn.Exec(query, data.UserId, data.HomestayId, data.DelState, data.Version, data.DeleteTime, data.Id)
	return err
}

func (m *defaultUserHomestayModel) Delete(id int64) error {
	query := `DELETE FROM ` + m.table + ` WHERE id = ?`
	_, err := m.conn.Exec(query, id)
	return err
}

func (m *defaultUserHomestayModel) CheckIfExists(ctx context.Context, userId, homestayId int64) (bool, error) {
	query := fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE user_id = ? AND homestay_id = ? AND del_state = 0", m.table)
	var count int
	err := m.conn.QueryRowCtx(ctx, &count, query, userId, homestayId)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (m *defaultUserHomestayModel) Favorite(userId, homestayId int64) error {
	query := `INSERT INTO ` + m.table + ` (user_id, homestay_id, del_state, version) VALUES (?, ?, 0, 1)`
	_, err := m.conn.Exec(query, userId, homestayId)
	return err
}

func (m *defaultUserHomestayModel) Unfavorite(userId, homestayId int64) error {
	query := `UPDATE ` + m.table + ` SET del_state = 1, delete_time = ? WHERE user_id = ? AND homestay_id = ? AND del_state = 0`
	_, err := m.conn.Exec(query, time.Now(), userId, homestayId)
	return err
}

// GetFavorites method with pagination
func (m *defaultUserHomestayModel) GetFavorites(userId int64, page int64, pageSize int64) ([]*UserHomestay, error) {
	offset := (page - 1) * pageSize
	query := `SELECT id, user_id, homestay_id, del_state, version, delete_time 
	          FROM ` + m.table + ` 
	          WHERE user_id = ? AND del_state = 0 
	          LIMIT ? OFFSET ?`
	var homestays []*UserHomestay
	err := m.conn.QueryRows(&homestays, query, userId, pageSize, offset)
	if err != nil {
		return nil, err
	}
	return homestays, nil
}

func (m *defaultUserHomestayModel) UpdateDelState(ctx context.Context, userId, homestayId int64, delState int) error {
	query := fmt.Sprintf("update %s set del_state = ? where user_id = ? and homestay_id = ? and del_state = 0", m.table)
	_, err := m.conn.ExecCtx(ctx, query, delState, userId, homestayId)
	return err
}
