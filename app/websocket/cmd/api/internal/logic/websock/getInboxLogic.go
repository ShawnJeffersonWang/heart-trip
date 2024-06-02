package websock

import (
	"context"
	"github.com/Masterminds/squirrel"
	"github.com/zeromicro/go-zero/core/mr"
	"golodge/app/websocket/model"
	"golodge/common/ctxdata"
	"golodge/common/globalkey"
	"sort"
	"strconv"

	"golodge/app/websocket/cmd/api/internal/svc"
	"golodge/app/websocket/cmd/api/internal/types"

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

// 用于按照 LastBrowsingTime 对 History 切片进行排序的自定义排序函数
func sortByCreateTime(messages []*types.Message) {
	sort.SliceStable(messages, func(i, j int) bool {
		return messages[i].CreateTime > messages[j].CreateTime
	})
}

func (l *GetInboxLogic) GetInbox(req *types.GetInboxReq) (*types.GetInboxResp, error) {
	// todo: add your logic here and delete this line
	userId := ctxdata.GetUidFromCtx(l.ctx)
	whereBuilder := l.svcCtx.MessageModel.SelectBuilder().Where(squirrel.Eq{
		"to_user_id": userId,
		"del_state":  globalkey.DelStateNo,
	})
	messages, _ := l.svcCtx.MessageModel.FindAll(l.ctx, whereBuilder, "create_time desc")
	//for _, message := range messages {
	//	fmt.Println("messages: ", *message)
	//}
	//fmt.Println("messages: ", *messages[0], *messages[1])
	var resp []*types.Message
	if len(messages) > 0 { // mapreduce example
		mr.MapReduceVoid(func(source chan<- interface{}) {
			for _, message := range messages {
				source <- message.Id
			}
		}, func(item interface{}, writer mr.Writer[*model.Message], cancel func(error)) {
			id := item.(int64)

			message, err := l.svcCtx.MessageModel.FindOne(l.ctx, id)
			if err != nil && err != model.ErrNotFound {
				logx.WithContext(l.ctx).Errorf("GetInbox 获取活动数据失败 id : %d ,err : %v", id, err)
				return
			}
			writer.Write(message)
		}, func(pipe <-chan *model.Message, cancel func(error)) {

			for message := range pipe {
				toId := strconv.Itoa(int(message.ToUserId))
				fromId := strconv.Itoa(int(message.FromUserId))
				tyMessage := types.Message{
					ToUserId:   toId,
					FromUserId: fromId,
					Content:    message.Content,
					Type:       "1",
					CreateTime: message.CreateTime.Unix(),
				}
				//_ = copier.Copy(&tyMessage, message)
				resp = append(resp, &tyMessage)
			}
		})
	}
	sortByCreateTime(resp)
	return &types.GetInboxResp{
		Messages: resp,
	}, nil
}
