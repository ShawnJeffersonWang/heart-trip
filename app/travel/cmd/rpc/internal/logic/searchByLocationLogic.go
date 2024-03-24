package logic

import (
	"context"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/jinzhu/copier"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/mr"
	"golodge/app/travel/model"
	"golodge/common/xerr"

	"golodge/app/travel/cmd/rpc/internal/svc"
	"golodge/app/travel/cmd/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type SearchByLocationLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSearchByLocationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SearchByLocationLogic {
	return &SearchByLocationLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SearchByLocationLogic) SearchByLocation(in *pb.SearchByLocationReq) (*pb.SearchByLocationResp, error) {
	// todo: add your logic here and delete this line
	if len(in.Location) == 0 {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DB_ERROR), "Search By Location err, len(Location) == 0")
	}
	// 模糊查询
	//whereBuilder := l.svcCtx.HomestayModel.SelectBuilder().Where("location LIKE ?", fmt.Sprint("%", in.Location, "%"))
	// 按地名或者详情查询
	whereBuilder := l.svcCtx.HomestayModel.SelectBuilder().Where(squirrel.Or{
		squirrel.Like{"location": fmt.Sprint("%", in.Location, "%")},
		squirrel.Like{"intro": fmt.Sprint("%", in.Location, "%")},
	})
	homestays, err := l.svcCtx.HomestayModel.FindAll(l.ctx, whereBuilder, "id desc")
	if err != nil {
		return nil, err
	}

	var resp []*pb.Homestay
	if len(homestays) > 0 {
		mr.MapReduceVoid(func(source chan<- any) {
			for _, homestay := range homestays {
				source <- homestay.Id
			}
		}, func(item any, writer mr.Writer[*model.Homestay], cancel func(error)) {
			id := item.(int64)
			homestay, err := l.svcCtx.HomestayModel.FindOne(l.ctx, id)
			if err != nil && err != model.ErrNotFound {
				logx.WithContext(l.ctx).Errorf("获取数据失败 id: %d err: %v", id, err)
				return
			}
			writer.Write(homestay)
		}, func(pipe <-chan *model.Homestay, cancel func(error)) {
			for homestay := range pipe {
				var tyHomestay pb.Homestay
				_ = copier.Copy(&tyHomestay, homestay)
				resp = append(resp, &tyHomestay)
			}
		})
	}

	return &pb.SearchByLocationResp{
		List: resp,
	}, nil
}
