package model

type Blog struct {
	ID     int64  `gorm:"primaryKey" json:"id"`
	UserID int64  `json:"user_id"`
	Name   string `json:"name"`
	Icon   string `json:"icon"`
	Liked  int64  `json:"liked"`
	IsLike bool   `json:"is_like"`
	// 其他字段...
}

type User struct {
	ID       int64  `gorm:"primaryKey" json:"id"`
	NickName string `json:"nick_name"`
	Icon     string `json:"icon"`
	// 其他字段...
}

type Follow struct {
	UserID       int64 `gorm:"column:user_id" json:"user_id"`
	FollowUserID int64 `gorm:"column:follow_user_id" json:"follow_user_id"`
	// 其他字段...
}
