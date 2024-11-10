package globalkey

import "time"

/**
redis key except "model cache key"  in here,
but "model cache key" in model
*/

const (
	// CacheUserTokenKey /** 用户登陆的token
	CacheUserTokenKey = "user_token:%d"
	CacheShopKey      = "cache:shop:"
	CacheNullTtl      = 5 * time.Minute
	LockShopKey       = "lock:shop:"
	BlogLikedKey      = "blog:liked:"
	FeedKey           = "feed:"
	ShopGeoKey        = "shop:geo:%d"
)
