package main

import (
	"flag"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
	"golodge/app/order/cmd/api/internal/config"
	"golodge/app/order/cmd/api/internal/handler"
	"golodge/app/order/cmd/api/internal/svc"
)

var configFile = flag.String("f", "etc/order.yaml", "the config file")

func main() {

	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	ctx := svc.NewServiceContext(c)
	server := rest.MustNewServer(c.RestConf, rest.WithCors("*"))
	defer server.Stop()

	handler.RegisterHandlers(server, ctx)

	server.Start()
}
