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
	whereBuilder := l.svcCtx.HistoryHomestayModel.SelectBuilder().Where(squirrel.Eq{
		"user_id": in.UserId,
	})
	historyIdList, _ := l.svcCtx.HistoryHomestayModel.FindAll(l.ctx, whereBuilder, "id desc")
	var resp []*pb.HistoryHomestay
	//fmt.Println("list: ", historyIdList)
	//for _, h := range historyIdList {
	//	fmt.Println("history: ", h)
	//}
	if len(historyIdList) > 0 {
		mr.MapReduceVoid(func(source chan<- any) {
			for _, historyHomestay := range historyIdList {
				source <- historyHomestay.HistoryId
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
