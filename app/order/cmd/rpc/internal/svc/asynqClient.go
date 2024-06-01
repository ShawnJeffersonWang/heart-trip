package svc

import (
	"github.com/hibiken/asynq"
	"golodge/app/order/cmd/rpc/internal/config"
)

// create asynq ws.
func newAsynqClient(c config.Config) *asynq.Client {
	return asynq.NewClient(asynq.RedisClientOpt{Addr: c.Redis.Host, Password: c.Redis.Pass})
}
