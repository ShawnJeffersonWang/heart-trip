package user

import (
	"heart-trip/common/result"
	"net/http"

	"heart-trip/app/usercenter/cmd/api/internal/logic/user"
	"heart-trip/app/usercenter/cmd/api/internal/svc"
	"heart-trip/app/usercenter/cmd/api/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
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
