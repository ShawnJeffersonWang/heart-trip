package config

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	rest.RestConf

	JwtAuth struct {
		AccessSecret string
	}
	DB struct {
		DataSource string
	}
	Cache cache.CacheConf

	UsercenterRpcConf zrpc.RpcClientConf
	TravelRpcConf     zrpc.RpcClientConf
}
