package user

import (
	"homestay/common/result"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"homestay/app/usercenter/cmd/api/internal/logic/user"
	"homestay/app/usercenter/cmd/api/internal/svc"
	"homestay/app/usercenter/cmd/api/internal/types"
)

func RemoveWishListHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.RemoveWishListReq
		if err := httpx.Parse(r, &req); err != nil {
			result.ParamErrorResult(r, w, err)
			return
		}

		l := user.NewRemoveWishListLogic(r.Context(), svcCtx)
		resp, err := l.RemoveWishList(&req)
		result.HttpResult(r, w, resp, err)
	}
}
