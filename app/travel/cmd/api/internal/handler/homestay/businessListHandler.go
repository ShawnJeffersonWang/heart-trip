package homestay

import (
	"net/http"

	"homestay/app/travel/cmd/api/internal/logic/homestay"
	"homestay/app/travel/cmd/api/internal/svc"
	"homestay/app/travel/cmd/api/internal/types"
	"homestay/common/result"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func BusinessListHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.BusinessListReq
		if err := httpx.Parse(r, &req); err != nil {
			result.ParamErrorResult(r, w, err)
			return
		}

		l := homestay.NewBusinessListLogic(r.Context(), ctx)
		resp, err := l.BusinessList(req)
		result.HttpResult(r, w, resp, err)
	}
}
