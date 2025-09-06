package mqs

import (
	"context"
	"fmt"
	"heart-trip/app/travel/cmd/mq/internal/svc"
)

type CommentSuccess struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCommentSuccess(ctx context.Context, svcCtx *svc.ServiceContext) *CommentSuccess {
	return &CommentSuccess{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CommentSuccess) Consume(ctx context.Context, key, val string) error {
	//logx.Infof("CommentSuccess key: %s val: %s", key, val)
	fmt.Printf("CommentSuccess key: %s val: %s", key, val)
	return nil
}
