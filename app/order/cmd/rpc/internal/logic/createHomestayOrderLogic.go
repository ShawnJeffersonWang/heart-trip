package logic

import (
	"context"
	"encoding/json"
	"fmt"
	"heart-trip/app/mqueue/cmd/job/jobtype"
	"strings"
	"time"

	"github.com/hibiken/asynq"

	"heart-trip/app/order/cmd/rpc/internal/svc"
	"heart-trip/app/order/cmd/rpc/pb"
	"heart-trip/app/order/model"
	"heart-trip/app/travel/cmd/rpc/travel"
	"heart-trip/common/tool"
	"heart-trip/common/uniqueid"
	"heart-trip/common/xerr"

	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
)

const CloseOrderTimeMinutes = 30 //defer close order time

type CreateHomestayOrderLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateHomestayOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateHomestayOrderLogic {
	return &CreateHomestayOrderLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// CreateHomestayOrder 会调用HomestayDetail, 而且是直接需要请求里面的resp.Homestay
func (l *CreateHomestayOrderLogic) CreateHomestayOrder(in *pb.CreateHomestayOrderReq) (*pb.CreateHomestayOrderResp, error) {

	//1、Create Order
	if in.LiveEndTime <= in.LiveStartTime {
		return nil, errors.Wrapf(xerr.NewErrMsg("Stay at least one night"), "Place an order at a B&B. The end time of your stay must be greater than the start time. in : %+v", in)
	}

	resp, err := l.svcCtx.TravelRpc.HomestayDetail(l.ctx, &travel.HomestayDetailReq{
		HomestayId: in.HomestayId,
	})
	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrMsg("Failed to query the record"), "Failed to query the record  rpc HomestayDetail fail , homestayId : %d , err : %v", in.HomestayId, err)
	}
	if resp == nil {
		return nil, errors.Wrapf(xerr.NewErrMsg("This record does not exist"), "This record does not exist , homestayId : %d ", in.HomestayId)
	}

	var cover string //Get the cover...
	// 看来我想的方法就是最佳实践: 多张图片使用一个字符串, 用','分割
	// 使用strings.Split分割然后用[0]取第一张图片
	// 这里左边的cover就是轮播图的第一张
	//if len(resp.Homestay.Cover) > 0 {
	//	cover = strings.Split(resp.Homestay.Cover, ",")[0]
	//}
	if len(resp.BannerUrls) > 0 {
		cover = strings.Split(resp.BannerUrls, ",")[0]
	}

	order := new(model.HomestayOrder)
	order.Sn = uniqueid.GenSn(uniqueid.SnPrefixHomestayOrder)
	order.UserId = in.UserId
	order.HomestayId = in.HomestayId
	order.Title = resp.Title
	order.Cover = cover
	order.Info = resp.RoomConfig
	order.HomestayPrice = resp.PriceBefore
	order.HomestayBusinessId = resp.HomestayBusinessId
	order.HomestayUserId = resp.HostId
	order.TradeState = model.HomestayOrderTradeStateWaitPay
	order.TradeCode = tool.Krand(8, tool.KC_RAND_KIND_ALL)
	order.Remark = in.Remark
	order.LiveStartDate = time.Unix(in.LiveStartTime, 0)
	order.LiveEndDate = time.Unix(in.LiveEndTime, 0)
	order.OrderTotalPrice = resp.PriceBefore

	liveDays := int64(order.LiveEndDate.Sub(order.LiveStartDate).Seconds() / 86400) //Stayed a few days in total

	order.HomestayTotalPrice = int64(resp.PriceBefore * liveDays) //Calculate the total price of the B&B
	fmt.Println(order.HomestayTotalPrice)

	_, err = l.svcCtx.HomestayOrderModel.Insert(l.ctx, nil, order)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DB_ERROR), "Order Database Exception order : %+v , err: %v", order, err)
	}

	//2、Delayed closing of order tasks.
	payload, err := json.Marshal(jobtype.DeferCloseHomestayOrderPayload{Sn: order.Sn})
	if err != nil {
		logx.WithContext(l.ctx).Errorf("create defer close order task json Marshal fail err :%+v , sn : %s", err, order.Sn)
	} else {
		_, err = l.svcCtx.AsynqClient.Enqueue(asynq.NewTask(jobtype.DeferCloseHomestayOrder, payload), asynq.ProcessIn(CloseOrderTimeMinutes*time.Minute))
		if err != nil {
			logx.WithContext(l.ctx).Errorf("create defer close order task insert queue fail err :%+v , sn : %s", err, order.Sn)
		}
	}

	return &pb.CreateHomestayOrderResp{
		Sn: order.Sn,
	}, nil
}
