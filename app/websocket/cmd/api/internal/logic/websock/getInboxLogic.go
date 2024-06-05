package websock

import (
	"context"
	"fmt"
	"golodge/app/websocket/cmd/api/internal/svc"
	"golodge/app/websocket/cmd/api/internal/types"
	"golodge/common/ctxdata"
	"strconv"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetInboxLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// get inbox messages
func NewGetInboxLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetInboxLogic {
	return &GetInboxLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetInboxLogic) GetInbox(req *types.GetInboxReq) (*types.GetInboxResp, error) {
	userId := ctxdata.GetUidFromCtx(l.ctx)
	messages, err := l.svcCtx.MessageModel.FindByUserId(l.ctx, userId)
	if err != nil {
		fmt.Println("查询失败=================")
		return nil, err
	}
	var responseMessages []types.Message
	for _, message := range messages {
		responseMessages = append(responseMessages, types.Message{
			FromUserId: strconv.FormatInt(message.FromUserId, 10),
			ToUserId:   strconv.FormatInt(message.ToUserId, 10),
			Content:    message.Content,
			CreateTime: message.CreateTime.Unix(),
		})
	}

	return &types.GetInboxResp{
		Messages: responseMessages,
	}, nil
}
