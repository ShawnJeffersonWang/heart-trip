package listen

import (
	"context"
	"heart-trip/app/travel/cmd/mq/internal/config"
	"heart-trip/app/travel/cmd/mq/internal/mqs"
	"heart-trip/app/travel/cmd/mq/internal/svc"

	"github.com/zeromicro/go-queue/kq"
	"github.com/zeromicro/go-zero/core/service"
)

func KqMqs(c config.Config, ctx context.Context, svcCtx *svc.ServiceContext) []service.Service {
	return []service.Service{
		kq.MustNewQueue(c.KqConsumerConf, mqs.NewCommentSuccess(ctx, svcCtx)),
	}
}
