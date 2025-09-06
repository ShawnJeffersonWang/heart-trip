package logic

import (
	"context"
	"fmt"
	"heart-trip/app/travel/cmd/rpc/internal/svc"
	"heart-trip/app/travel/cmd/rpc/pb"
	"heart-trip/common/globalkey"
	"os"
	"strconv"

	"github.com/zeromicro/go-zero/core/logx"
)

type SeckillVoucherLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSeckillVoucherLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SeckillVoucherLogic {
	return &SeckillVoucherLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SeckillVoucherLogic) SeckillVoucher(req *pb.SeckillVoucherRequest) (*pb.SeckillVoucherResponse, error) {
	//userId := ctxdata.GetUidFromCtx(l.ctx)
	voucherId := req.VoucherId

	//orderId := l.svcCtx.RedisClient.Incr(l.ctx, "order:id").Val() // 使用 Redis ID 生成器
	// redis上下文只注入一次，防止每次请求都建立 redis 连接，释放 redis 连接，产生开销
	//redisIdWorker := tool.NewRedisIdWorker(redis.NewClient(&redis.Options{
	//	Addr:     l.svcCtx.Config.Cache[0].Host,
	//	Password: l.svcCtx.Config.Cache[0].Pass,
	//}))
	orderId, _ := l.svcCtx.RedisIdWorker.NextID("order")
	// 执行 Lua 脚本
	script, err := os.ReadFile(globalkey.SeckillScriptPath)
	if err != nil {
		return nil, fmt.Errorf("读取脚本文件失败: %v", err)
	}

	result, err := l.svcCtx.RedisClient.Eval(l.ctx,
		string(script), nil,
		strconv.FormatInt(voucherId, 10),
		req.UserId,
		strconv.FormatInt(orderId, 10)).Result()

	if err != nil {
		l.Logger.Error("执行 Lua 脚本失败: ", err)
		return &pb.SeckillVoucherResponse{
			Code:    500,
			Message: "内部错误",
		}, nil
	}

	// 处理 Lua 脚本返回结果
	switch res := result.(type) {
	case int64:
		if res != 0 {
			return &pb.SeckillVoucherResponse{
				Code:    res,
				Message: getMessageByCode(res),
			}, nil
		}
	default:
		return &pb.SeckillVoucherResponse{
			Code:    500,
			Message: "内部错误",
		}, nil
	}

	return &pb.SeckillVoucherResponse{
		Code:    200,
		Message: "下单成功",
		OrderId: orderId,
	}, nil
}

// getMessageByCode 根据返回码获取消息
func getMessageByCode(code int64) string {
	switch code {
	case 1:
		return "库存不足"
	case 2:
		return "不能重复下单"
	default:
		return "未知错误"
	}
}
