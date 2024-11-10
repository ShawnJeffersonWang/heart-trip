package model

import (
	"context"
	"errors"
	"gorm.io/gorm"
)

type ShopModel interface {
	FindOne(ctx context.Context, id int64) (*Homestay, error)
}

// 实现 ShopModel 接口
type shopModel struct {
	db *gorm.DB
}

func NewShopModel(db *gorm.DB) ShopModel {
	return &shopModel{db: db}
}

func (m *shopModel) FindOne(ctx context.Context, id int64) (*Homestay, error) {
	var shop Homestay
	result := m.db.WithContext(ctx).First(&shop, id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, errors.New("shop not found")
	}
	if result.Error != nil {
		return nil, result.Error
	}
	return &shop, nil
}
