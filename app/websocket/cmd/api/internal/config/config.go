package config

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	rest.RestConf
	//Path    string `json:",default=."`
	JwtAuth struct {
		AccessSecret string
	}
	//OSS_ACCESS_KEY_ID     string `json:",env=OSS_ACCESS_KEY_ID"`
	//OSS_ACCESS_KEY_SECRET string `json:",env=OSS_ACCESS_KEY_SECRET"`
	UsercenterRpcConf zrpc.RpcClientConf
	TravelRpcConf     zrpc.RpcClientConf
	DB                struct {
		DataSource string
	}
	Cache cache.CacheConf
}
