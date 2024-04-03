package main

import (
	"flag"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
	"golodge/app/travel/cmd/mq/internal/config"
	"golodge/app/travel/cmd/mq/internal/listen"
)

var configFile = flag.String("f", "etc/travel.yaml", "the config file")

func main() {
	flag.Parse()
	var c config.Config
	conf.MustLoad(*configFile, &c)
	serviceGroup := service.NewServiceGroup()
	defer serviceGroup.Stop()
	for _, mq := range listen.Mqs(c) {
		serviceGroup.Add(mq)
	}
	serviceGroup.Start()
}
