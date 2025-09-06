package homestay

import (
	"net/http"

	"heart-trip/app/travel/cmd/api/internal/logic/homestay"
	"heart-trip/app/travel/cmd/api/internal/svc"
	"heart-trip/app/travel/cmd/api/internal/types"
	"heart-trip/common/result"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func HomestayListHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.HomestayListReq
		if err := httpx.Parse(r, &req); err != nil {
			result.ParamErrorResult(r, w, err)
			return
		}

		l := homestay.NewHomestayListLogic(r.Context(), ctx)
		resp, err := l.HomestayList(req)
		result.HttpResult(r, w, resp, err)
	}
}
