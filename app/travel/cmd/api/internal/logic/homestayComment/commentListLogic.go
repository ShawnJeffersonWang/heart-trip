package homestayComment

import (
	"context"
	"github.com/Masterminds/squirrel"
	"github.com/jinzhu/copier"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/mr"
	"golodge/app/travel/model"
	"golodge/common/xerr"

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

			for homestay := range pipe {
				var tyHomestay types.HomestayComment
				_ = copier.Copy(&tyHomestay, homestay)

				resp = append(resp, tyHomestay)
			}
		})
	}

	return &types.CommentListResp{
		List: resp,
	}, nil
}
