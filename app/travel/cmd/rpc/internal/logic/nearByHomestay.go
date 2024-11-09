// shop/internal/logic/queryshopbytypelogic.go

package logic

import (
	"context"
	"fmt"
	"golodge/app/travel/cmd/rpc/internal/svc"
	"golodge/app/travel/cmd/rpc/pb"
	"golodge/common/globalkey"
	"strconv"
	"strings"

	"github.com/go-redis/redis/v8"
	"github.com/zeromicro/go-zero/core/logx"
)

type QueryShopByTypeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewQueryShopByTypeLogic(ctx context.Context, svcCtx *svc.ServiceContext) QueryShopByTypeLogic {
	return QueryShopByTypeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *QueryShopByTypeLogic) QueryShopByType(req *pb.QueryShopByTypeRequest) (*pb.QueryShopByTypeResponse, error) {
	ctx := context.Background()
	typeId := req.TypeId
	current := req.Current
	x := req.X
	y := req.Y

	if current <= 0 {
		current = 1 // 保证当前页数至少为1
	}

	if x == 0 && y == 0 {
		// 不需要根据坐标查询，直接从数据库分页查询
		var shops []*pb.Homestay
		offset := (current - 1) * globalkey.DefaultPageSize
		limit := globalkey.DefaultPageSize

		result := l.svcCtx.DB.Where("type_id = ?", typeId).
			Offset(int(offset)).
			Limit(limit).
			Find(&shops)

		if result.Error != nil {
			logx.Errorw("数据库查询失败", logx.LogField{Key: "error", Value: result.Error})
			return &pb.QueryShopByTypeResponse{
				Code: 500,
				Msg:  "数据库查询失败: " + result.Error.Error(),
				Data: nil,
			}, nil
		}

		return &pb.QueryShopByTypeResponse{
			Code: 200,
			Msg:  "success",
			Data: shops,
		}, nil
	}

	// 需要根据坐标查询
	offset := (current - 1) * globalkey.DefaultPageSize
	end := current * globalkey.DefaultPageSize

	key := fmt.Sprintf(globalkey.ShopGeoKey, typeId)
	// 使用 Redis GEOSEARCH
	geoArgs := &redis.GeoSearchLocationQuery{
		GeoSearchQuery: redis.GeoSearchQuery{
			Longitude:  x,
			Latitude:   y,
			Radius:     5000, // 半径 5000 米
			RadiusUnit: "m",
			Sort:       "ASC",
			Count:      int(end),
		},
		WithDist: true,
	}

	geoResults, err := l.svcCtx.RedisClient.GeoSearchLocation(ctx, key, geoArgs).Result()
	if err != nil {
		logx.Errorw("Redis GEOSEARCH 失败", logx.LogField{Key: "error", Value: err})
		return &pb.QueryShopByTypeResponse{
			Code: 500,
			Msg:  "Redis GEOSEARCH 失败: " + err.Error(),
			Data: nil,
		}, nil
	}

	if len(geoResults) == 0 {
		// 没有结果
		return &pb.QueryShopByTypeResponse{
			Code: 200,
			Msg:  "success",
			Data: []*pb.Homestay{},
		}, nil
	}

	if len(geoResults) <= int(offset) {
		// 没有下一页数据
		return &pb.QueryShopByTypeResponse{
			Code: 200,
			Msg:  "success",
			Data: []*pb.Homestay{},
		}, nil
	}

	// 截取当前页需要的数据
	slicedResults := geoResults[offset:]
	if len(slicedResults) > globalkey.DefaultPageSize {
		slicedResults = slicedResults[:globalkey.DefaultPageSize]
	}

	// 收集 Shop ID 和距离
	ids := make([]int64, 0, len(slicedResults))
	distanceMap := make(map[string]float64, len(slicedResults))
	for _, loc := range slicedResults {
		shopIDStr := loc.Name
		shopID, err := strconv.ParseInt(shopIDStr, 10, 64)
		if err != nil {
			// 跳过无效的 Shop ID
			logx.Errorw("无效的 Shop ID", logx.LogField{Key: "shopIDStr", Value: shopIDStr}, logx.LogField{Key: "error", Value: err})
			continue
		}
		ids = append(ids, shopID)
		distanceMap[shopIDStr] = loc.Dist
	}

	if len(ids) == 0 {
		// 没有有效的 Shop ID
		return &pb.QueryShopByTypeResponse{
			Code: 200,
			Msg:  "success",
			Data: []*pb.Homestay{},
		}, nil
	}

	// 查询数据库中对应的商铺信息，并按 Redis 返回的顺序排序
	var shops []*pb.Homestay
	// 使用 GORM 的 Order 方法进行 FIELD 排序
	idInterfaces := make([]interface{}, len(ids))
	for i, id := range ids {
		idInterfaces[i] = id
	}

	idListStr := make([]string, len(ids))
	for i, id := range ids {
		idListStr[i] = strconv.FormatInt(id, 10)
	}
	orderBy := fmt.Sprintf("FIELD(id, %s)", strings.Join(idListStr, ","))

	result := l.svcCtx.DB.Where("id IN ?", ids).
		Order(orderBy).
		Find(&shops)

	if result.Error != nil {
		logx.Errorw("数据库查询失败", logx.LogField{Key: "error", Value: result.Error})
		return &pb.QueryShopByTypeResponse{
			Code: 500,
			Msg:  "数据库查询失败: " + result.Error.Error(),
			Data: nil,
		}, nil
	}

	// 将距离信息赋值到 Shop 结构体
	for i := range shops {
		shopIDStr := strconv.FormatInt(shops[i].Id, 10)
		if dist, ok := distanceMap[shopIDStr]; ok {
			shops[i].Distance = dist
		}
	}

	return &pb.QueryShopByTypeResponse{
		Code: 200,
		Msg:  "success",
		Data: shops,
	}, nil
}
