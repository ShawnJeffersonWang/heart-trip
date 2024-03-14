package user

import (
	"homestay/common/result"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"homestay/app/usercenter/cmd/api/internal/logic/user"
	"homestay/app/usercenter/cmd/api/internal/svc"
	"homestay/app/usercenter/cmd/api/internal/types"
)

func WishListHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.WishListReq
		if err := httpx.Parse(r, &req); err != nil {
			result.ParamErrorResult(r, w, err)
			return
		}

		l := user.NewWishListLogic(r.Context(), svcCtx)
		resp, err := l.WishList(&req)
		result.HttpResult(r, w, resp, err)
	}
}
