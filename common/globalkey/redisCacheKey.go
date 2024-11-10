package globalkey

/**
redis key except "model cache key"  in here,
but "model cache key" in model
*/

const (
	// CacheUserTokenKey /** 用户登陆的token
	CacheUserTokenKey = "user_token:%d"
	BlogLikedKey      = "blog:liked:"
	FeedKey           = "feed:"
	ShopGeoKey        = "shop:geo:%d"
)
