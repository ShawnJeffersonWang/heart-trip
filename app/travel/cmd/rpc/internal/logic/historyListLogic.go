package logic

import (
	"context"
	"github.com/Masterminds/squirrel"
	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/mr"
	"golodge/app/travel/model"

	"golodge/app/travel/cmd/rpc/internal/svc"
	"golodge/app/travel/cmd/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type HistoryListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewHistoryListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *HistoryListLogic {
	return &HistoryListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *HistoryListLogic) HistoryList(in *pb.HistoryListReq) (*pb.HistoryListResp, error) {
	// todo: add your logic here and delete this line
	whereBuilder := l.svcCtx.UserHistoryModel.SelectBuilder().Where(squirrel.Eq{
		"user_id": in.UserId,
	})
	// 按历史记录时间倒叙排列
	userHistories, _ := l.svcCtx.UserHistoryModel.FindAll(l.ctx, whereBuilder, "")
	var resp []*pb.HistoryHomestay
	//fmt.Println("list: ", userHistories)
	//for _, h := range userHistories {
	//	fmt.Println("history: ", h)
	//}
	if len(userHistories) > 0 {
		mr.MapReduceVoid(func(source chan<- any) {
			for _, userHistory := range userHistories {
				source <- userHistory.HistoryId
			}
		}, func(item any, writer mr.Writer[*model.History], cancel func(error)) {
			id := item.(int64)
			history, err := l.svcCtx.HistoryModel.FindOne(l.ctx, id)
			if err != nil && err != model.ErrNotFound {
				logx.WithContext(l.ctx).Errorf("获取活动数据失败 id: %d err: %v", id, err)
				return
			}
			writer.Write(history)
		}, func(pipe <-chan *model.History, cancel func(error)) {
			for history := range pipe {
				var tyHistory pb.HistoryHomestay
				_ = copier.Copy(&tyHistory, history)
				resp = append(resp, &tyHistory)
			}
		})
	}

	return &pb.HistoryListResp{
		HistoryList: resp,
	}, nil
}
