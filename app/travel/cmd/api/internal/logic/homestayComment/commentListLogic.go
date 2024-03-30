package homestayComment

import (
	"context"
	"github.com/Masterminds/squirrel"
	"github.com/jinzhu/copier"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/mr"
	"golodge/app/travel/model"
	"golodge/common/xerr"
	"strconv"

	"golodge/app/travel/cmd/api/internal/svc"
	"golodge/app/travel/cmd/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CommentListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCommentListLogic(ctx context.Context, svcCtx *svc.ServiceContext) CommentListLogic {
	return CommentListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CommentListLogic) CommentList(req types.CommentListReq) (*types.CommentListResp, error) {
	// todo: add your logic here and delete this line

	whereBuilder := l.svcCtx.HomestayCommentModel.SelectBuilder().Where(squirrel.Eq{
		"homestay_id": req.HomestayId,
	})
	homestayCommentList, err := l.svcCtx.HomestayCommentModel.FindPageListByPage(l.ctx, whereBuilder, req.Page, req.PageSize, "id desc")
	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DB_ERROR), "get activity homestay id set fail rowType: %s ,err : %v", model.HomestayActivityPreferredType, err)
	}

	var resp []types.HomestayComment
	var starSum float64 = 0
	var TidyRatingSum float64 = 0
	var TrafficRatingSum float64 = 0
	var SecurityRatingSum float64 = 0
	var FoodRatingSum float64 = 0
	var CostRatingSum float64 = 0

	if len(homestayCommentList) > 0 { // mapreduce example
		mr.MapReduceVoid(func(source chan<- interface{}) {
			for _, homestayComment := range homestayCommentList {
				source <- homestayComment.Id
			}
		}, func(item interface{}, writer mr.Writer[*model.HomestayComment], cancel func(error)) {
			id := item.(int64)

			homestayComment, err := l.svcCtx.HomestayCommentModel.FindOne(l.ctx, id)
			if err != nil && err != model.ErrNotFound {
				logx.WithContext(l.ctx).Errorf("ActivityHomestayListLogic ActivityHomestayList 获取活动数据失败 id : %d ,err : %v", id, err)
				return
			}
			writer.Write(homestayComment)
		}, func(pipe <-chan *model.HomestayComment, cancel func(error)) {

			for homestayComment := range pipe {
				var comment types.HomestayComment
				_ = copier.Copy(&comment, homestayComment)
				star, _ := strconv.ParseFloat(comment.Star, 64)
				tidyRating, _ := strconv.ParseFloat(comment.TidyRating, 64)
				trafficRating, _ := strconv.ParseFloat(comment.TrafficRating, 64)
				securityRating, _ := strconv.ParseFloat(comment.SecurityRating, 64)
				foodRating, _ := strconv.ParseFloat(comment.FoodRating, 64)
				costRating, _ := strconv.ParseFloat(comment.CostRating, 64)
				starSum += star
				TidyRatingSum += tidyRating
				TrafficRatingSum += trafficRating
				SecurityRatingSum += securityRating
				FoodRatingSum += foodRating
				CostRatingSum += costRating
				resp = append(resp, comment)
			}
		})
	}
	count := len(homestayCommentList)
	cnt := float64(count)
	return &types.CommentListResp{
		HomestayId:     req.HomestayId,
		Star:           strconv.FormatFloat(starSum/cnt, 'f', 2, 64),
		TidyRating:     strconv.FormatFloat(TidyRatingSum/cnt, 'f', 2, 64),
		TrafficRating:  strconv.FormatFloat(TrafficRatingSum/cnt, 'f', 2, 64),
		SecurityRating: strconv.FormatFloat(SecurityRatingSum/cnt, 'f', 2, 64),
		FoodRating:     strconv.FormatFloat(FoodRatingSum/cnt, 'f', 2, 64),
		CostRating:     strconv.FormatFloat(CostRatingSum/cnt, 'f', 2, 64),
		CommentCount:   int64(count),
		List:           resp,
	}, nil
}
