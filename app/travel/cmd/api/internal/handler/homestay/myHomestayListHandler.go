package homestay

import (
	"heart-trip/common/result"
	"net/http"

	"heart-trip/app/travel/cmd/api/internal/logic/homestay"
	"heart-trip/app/travel/cmd/api/internal/svc"
	"heart-trip/app/travel/cmd/api/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

// my homestay room list
func MyHomestayListHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.MyHomestayListReq
		if err := httpx.Parse(r, &req); err != nil {
			result.ParamErrorResult(r, w, err)
			return
		}

		l := homestay.NewMyHomestayListLogic(r.Context(), svcCtx)
		resp, err := l.MyHomestayList(&req)
		result.HttpResult(r, w, resp, err)
	}
}
