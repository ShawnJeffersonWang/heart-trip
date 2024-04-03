package logic

import (
	"context"
	"github.com/Masterminds/squirrel"
	"github.com/zeromicro/go-zero/core/mr"
	"golodge/app/travel/model"
	"sort"

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

// 用于按照 LastBrowsingTime 对 History 切片进行排序的自定义排序函数
func sortByLastBrowsingTime(histories []*pb.History) {
	sort.SliceStable(histories, func(i, j int) bool {
		return histories[i].LastBrowsingTime > histories[j].LastBrowsingTime
	})
}

func (l *HistoryListLogic) HistoryList(in *pb.HistoryListReq) (*pb.HistoryListResp, error) {
	// todo: add your logic here and delete this line
	whereBuilder := l.svcCtx.UserHistoryModel.SelectBuilder().Where(squirrel.Eq{
		"user_id": in.UserId,
	})
	// 按历史记录时间倒叙排列, bug: 这里不能按照更新时间，因为UserHistory这个表压根没有update_time这个字段
	userHistories, _ := l.svcCtx.UserHistoryModel.FindAll(l.ctx, whereBuilder, "")
	var resp []*pb.History
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
				//var tyHistory pb.History
				// 爽, 可以用Fill all fields一键填充所有字段
				// 不用copier.Copy, 手动映射lastBrowsingTime
				tyHistory := pb.History{
					Id:                 history.Id,
					Title:              history.Title,
					Cover:              history.Cover,
					Intro:              history.Intro,
					Location:           history.Location,
					HomestayBusinessId: history.HomestayBusinessId,
					UserId:             history.UserId,
					RowState:           history.RowState,
					RatingStars:        history.RatingStars,
					PriceBefore:        history.PriceBefore,
					PriceAfter:         history.PriceAfter,
					LastBrowsingTime:   history.LastBrowsingTime.Unix(),
				}
				//_ = copier.Copy(&tyHistory, history)
				resp = append(resp, &tyHistory)
			}
		})
	}

	sortByLastBrowsingTime(resp)
	return &pb.HistoryListResp{
		HistoryList: resp,
	}, nil
}
