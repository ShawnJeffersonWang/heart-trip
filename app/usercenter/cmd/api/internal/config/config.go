package config

import (
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
	WxMiniConf        WxMiniConf
	UsercenterRpcConf zrpc.RpcClientConf
	TravelRpcConf     zrpc.RpcClientConf
}
