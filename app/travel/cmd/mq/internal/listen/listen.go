package listen

import (
	"context"
	"github.com/zeromicro/go-zero/core/service"
	"golodge/app/travel/cmd/mq/internal/config"
	"golodge/app/travel/cmd/mq/internal/svc"
)

// back to all consumers
func Mqs(c config.Config) []service.Service {

	svcContext := svc.NewServiceContext(c)
	ctx := context.Background()

	var services []service.Service

	//kq ：pub sub
	services = append(services, KqMqs(c, ctx, svcContext)...)

	return services
}
