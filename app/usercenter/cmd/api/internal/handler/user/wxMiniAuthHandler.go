package user

import (
	"net/http"

	"heart-trip/app/usercenter/cmd/api/internal/logic/user"
	"heart-trip/app/usercenter/cmd/api/internal/svc"
	"heart-trip/app/usercenter/cmd/api/internal/types"
	"heart-trip/common/result"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func WxMiniAuthHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.WXMiniAuthReq
		if err := httpx.Parse(r, &req); err != nil {
			result.ParamErrorResult(r, w, err)
			return
		}

		l := user.NewWxMiniAuthLogic(r.Context(), ctx)
		resp, err := l.WxMiniAuth(req)
		result.HttpResult(r, w, resp, err)
	}
}
