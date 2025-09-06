package listen

import (
	"context"
	"golodge/app/order/cmd/mq/internal/config"
	"golodge/app/order/cmd/mq/internal/svc"

	"github.com/zeromicro/go-zero/core/service"
)

// Mqs back to all consumers
func Mqs(c config.Config) []service.Service {

	svcContext := svc.NewServiceContext(c)
	ctx := context.Background()

	var services []service.Service

	//kq ：pub sub
	services = append(services, KqMqs(c, ctx, svcContext)...)

	return services
}
