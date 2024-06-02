package main

import (
	"flag"
	"fmt"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
	"golodge/app/websocket/cmd/api/internal/config"
	"golodge/app/websocket/cmd/api/internal/handler"
	"golodge/app/websocket/cmd/api/internal/logic/process"
	"golodge/app/websocket/cmd/api/internal/svc"
)

var configFile = flag.String("f", "etc/websocket.yaml", "the config file")

func main() {
	flag.Parse()
	var c config.Config
	conf.MustLoad(*configFile, &c, conf.UseEnv())

	ctx := svc.NewServiceContext(c)
	server := rest.MustNewServer(c.RestConf, rest.WithCors("*"))
	defer server.Stop()

	hub := process.NewHub()
	go hub.Run()
	handler.RegisterHandlers(server, ctx, hub)

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
