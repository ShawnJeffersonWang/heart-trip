package logic

import (
	"context"
	"golodge/app/travel/model"

	"gorm.io/gorm"
)

// FollowService 实现 IFollowService 接口
type FollowService struct {
	db *gorm.DB
}

// NewFollowService 创建一个新的 FollowService 实例
func NewFollowService(db *gorm.DB) *FollowService {
	return &FollowService{
		db: db,
	}
}

// QueryFollowsByFollowUserID 根据关注用户ID查询所有关注者
func (s *FollowService) QueryFollowsByFollowUserID(ctx context.Context, followUserID int64, follows *[]model.Follow) error {
	// 使用上下文并行控制，例如超时、取消等
	return s.db.WithContext(ctx).Where("follow_user_id = ?", followUserID).Find(follows).Error
}
