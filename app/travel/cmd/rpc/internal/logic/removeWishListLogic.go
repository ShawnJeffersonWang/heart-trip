package logic

import (
	"context"
	"github.com/zeromicro/go-zero/core/logx"
	"golodge/app/travel/cmd/rpc/internal/svc"
	"golodge/app/travel/cmd/rpc/pb"
)

type RemoveWishListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRemoveWishListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RemoveWishListLogic {
	return &RemoveWishListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *RemoveWishListLogic) RemoveWishList(in *pb.RemoveWishListReq) (*pb.RemoveWishListResp, error) {
	//userHomestay, err := l.svcCtx.UserHomestayModel.FindOneByUserIdAndHomestayId(l.ctx, in.UserId, in.HomestayId)
	//err = l.svcCtx.UserHistoryModel.Trans(l.ctx, func(ctx context.Context, session sqlx.Session) error {
	//	err := l.svcCtx.UserHomestayModel.DeleteSoft(ctx, session, userHomestay)
	//	if err != nil {
	//		return err
	//	}
	//	return nil
	//})
	//if err != nil {
	//	return nil, err
	//}
	//return &pb.RemoveWishListResp{}, nil

	// 检查用户是否已经收藏了该民宿
	exists, err := l.svcCtx.UserHomestayModel.CheckIfExists(l.ctx, in.UserId, in.HomestayId)
	if err != nil {
		return nil, err
	}
	if !exists {
		return &pb.RemoveWishListResp{
			Success: false,
		}, nil
	}

	// 更新 del_state 字段以取消收藏
	err = l.svcCtx.UserHomestayModel.UpdateDelState(l.ctx, in.UserId, in.HomestayId, 1)
	if err != nil {
		return nil, err
	}
	//homestay, err := l.svcCtx.HomestayModel.FindOne(l.ctx, in.HomestayId)
	//homestay.I
	//l.svcCtx.HomestayModel.Update(l.ctx,in.)
	//updateHomestayLogic:=NewUpdateHomestayLogic(l.ctx,l.svcCtx)
	//updateHomestayLogic.UpdateHomestay(homestay.)
	//if err := l.svcCtx.HomestayModel.Trans(l.ctx, func(ctx context.Context, session sqlx.Session) error {
	//	homestay := model.Homestay{
	//		Id:           in.HomestayId,
	//		DeleteTime:   time.Now(),
	//		DelState:     0,
	//		Version:      0,
	//		Title:        homestay.Title,
	//		BannerUrls:   homestay.BannerUrls,
	//		TitleTags:    homestay.TitleTags,
	//		Latitude:     homestay.Latitude,
	//		Longitude:    homestay.Longitude,
	//		Location:     homestay.Location,
	//		Facilities:   homestay.Facilities,
	//		Cover:        homestay.Cover,
	//		Area:         homestay.Area,
	//		RoomConfig:   homestay.RoomConfig,
	//		CleanVideo:   homestay.CleanVideo,
	//		HostId:       homestay.HostId,
	//		HostAvatar:   homestay.HostAvatar,
	//		HostNickname: homestay.HostNickname,
	//		RowState:     homestay.RowState,
	//		PriceBefore:  homestay.PriceBefore,
	//		PriceAfter:   homestay.PriceAfter,
	//	}
	//
	//	res, err := l.svcCtx.HomestayModel.Update(ctx, session, &homestay)
	//	if err != nil {
	//		log.Println("travel.UpdateHomestay.err: ", err)
	//		return err
	//	}
	//
	//	// bug: 这里可以直接通过homestay获取新加入的Id, 而不用再去查询, 多此一举, 因为是通过&homestay获取的
	//	// 就像C语言的传入传出参数一样, 更正: 虽然类似传出参数, 但是不能获取自增的id
	//	// 之前加入的无法在列表中拿到是因为之加入了Homestay表, 而没有加入HomestayActivity表
	//	// 困扰了n天的bug: 加入数据库后自增的id使用sql.Result接口中的LastInsertId()方法获取
	//	dataId, _ := res.LastInsertId()
	//	homestayActivity := model.HomestayActivity{
	//		DelState:  0,
	//		RowType:   "preferredHomestay",
	//		DataId:    dataId,
	//		RowStatus: 1,
	//		Version:   0,
	//	}
	//
	//	_, err = l.svcCtx.HomestayActivityModel.Update(ctx, session, &homestayActivity)
	//	if err != nil {
	//		return err
	//	}
	//
	//	return nil
	//}); err != nil {
	//	return nil, err
	//}

	return &pb.RemoveWishListResp{
		Success: true,
	}, nil
}
