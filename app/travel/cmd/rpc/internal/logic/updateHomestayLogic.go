package logic

import (
	"context"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"golodge/app/travel/cmd/rpc/travel"
	"golodge/app/travel/model"
	"log"
	"strings"
	"time"

	"golodge/app/travel/cmd/rpc/internal/svc"
	"golodge/app/travel/cmd/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateHomestayLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateHomestayLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateHomestayLogic {
	return &UpdateHomestayLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateHomestayLogic) UpdateHomestay(in *pb.UpdateHomestayReq) (*pb.UpdateHomestayResp, error) {
	// todo: add your logic here and delete this line
	//_, err := l.svcCtx.HomestayModel.FindOne(l.ctx, in.Homestay.Id)
	//if err == nil {
	//	return nil, errors.Wrapf(ErrHomestayAlreadyAdded,
	//		"homestay has been added in homestayList homestayId:%d,err:%v", in.Homestay.Id, err)
	//}
	homestay, err := l.svcCtx.HomestayModel.FindOne(l.ctx, in.HomestayId)
	if err != nil {
		return nil, err
	}
	if in.HostId != homestay.HostId {
		return nil, errors.New("No authorization")
	}
	splitUrls := strings.Split(in.BannerUrls, ",")
	firstUrl := splitUrls[0]
	if err := l.svcCtx.HomestayModel.Trans(l.ctx, func(ctx context.Context, session sqlx.Session) error {
		homestay := model.Homestay{
			Id:           in.HomestayId,
			DeleteTime:   time.Now(),
			DelState:     0,
			Version:      0,
			Title:        in.Title,
			BannerUrls:   in.BannerUrls,
			TitleTags:    in.TitleTags,
			Latitude:     in.Latitude,
			Longitude:    in.Longitude,
			Location:     in.Location,
			Facilities:   in.Facilities,
			Cover:        firstUrl,
			Area:         in.Area,
			RoomConfig:   in.RoomConfig,
			CleanVideo:   in.CleanVideo,
			HostId:       in.HostId,
			HostAvatar:   in.HostAvatar,
			HostNickname: in.HostNickname,
			RowState:     in.RowState,
			PriceBefore:  in.PriceBefore,
			PriceAfter:   in.PriceAfter,
		}

		res, err := l.svcCtx.HomestayModel.Update(ctx, session, &homestay)
		if err != nil {
			log.Println("travel.UpdateHomestay.err: ", err)
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

		_, err = l.svcCtx.HomestayActivityModel.Update(ctx, session, &homestayActivity)
		if err != nil {
			return err
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return &travel.UpdateHomestayResp{}, nil
}
