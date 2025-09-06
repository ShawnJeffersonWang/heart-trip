package main

import (
	"flag"

	"heart-trip/app/order/cmd/api/internal/config"
	"heart-trip/app/order/cmd/api/internal/handler"
	"heart-trip/app/order/cmd/api/internal/svc"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
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
