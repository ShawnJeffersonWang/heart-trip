package model

import (
	"errors"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var ErrNotFound = sqlx.ErrNotFound
var ErrNoRowsUpdate = errors.New("update db no rows change")

var UserAuthTypeSystem = "system"  //平台内部
var UserAuthTypeSmallWX = "wxMini" //微信小程序
