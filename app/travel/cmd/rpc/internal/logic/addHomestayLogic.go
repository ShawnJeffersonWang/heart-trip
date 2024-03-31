package logic

import (
	"context"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"golodge/app/travel/cmd/rpc/internal/svc"
	"golodge/app/travel/cmd/rpc/pb"
	"golodge/app/travel/cmd/rpc/travel"
	"golodge/app/travel/model"
	"strings"
)

type AddHomestayLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAddHomestayLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddHomestayLogic {
	return &AddHomestayLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *AddHomestayLogic) AddHomestay(in *pb.AddHomestayReq) (*pb.AddHomestayResp, error) {
	// todo: add your logic here and delete this line
	//_, err := l.svcCtx.HomestayModel.FindOne(l.ctx, in.Homestay.Id)
	//if err == nil {
	//	return nil, errors.Wrapf(ErrHomestayAlreadyAdded,
	//		"homestay has been added in homestayList homestayId:%d,err:%v", in.Homestay.Id, err)
	//}
	splitUrls := strings.Split(in.BannerUrls, ",")
	firstUrl := splitUrls[0]
	if err := l.svcCtx.HomestayModel.Trans(l.ctx, func(ctx context.Context, session sqlx.Session) error {
		homestay := model.Homestay{
			DelState:     0,
			Version:      0,
			Title:        in.Title,
			BannerUrls:   in.BannerUrls,
			TitleTags:    in.TitleTags,
			Latitude:     in.Latitude,
			Longitude:    in.Longitude,
			Facilities:   in.Facilities,
			Cover:        firstUrl,
			Area:         in.Area,
			RoomConfig:   in.RoomConfig,
			CleanVideo:   in.CleanVideo,
			HostId:       in.HostId,
			HostAvatar:   in.HostAvatar,
			HostNickname: in.HostNickname,
			RowState:     1,
			PriceBefore:  in.PriceBefore,
			PriceAfter:   in.PriceAfter,
		}

		res, err := l.svcCtx.HomestayModel.Insert(ctx, session, &homestay)
		if err != nil {
			return err
		}

		// bug: 这里可以直接通过homestay获取新加入的Id, 而不用再去查询, 多此一举, 因为是通过&homestay获取的
		// 就像C语言的传入传出参数一样, 更正: 虽然类似传出参数, 但是不能获取自增的id
		// 之前加入的无法在列表中拿到是因为之加入了Homestay表, 而没有加入HomestayActivity表
		// 困扰了n天的bug: 加入数据库后自增的id使用sql.Result接口中的LastInsertId()方法获取
		dataId, _ := res.LastInsertId()
		homestayActivity := model.HomestayActivity{
			DelState:  0,
			RowType:   "preferredHomestay",
			DataId:    dataId,
			RowStatus: 1,
			Version:   0,
		}

		_, err = l.svcCtx.HomestayActivityModel.Insert(ctx, session, &homestayActivity)
		if err != nil {
			return err
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return &travel.AddHomestayResp{}, nil
}
