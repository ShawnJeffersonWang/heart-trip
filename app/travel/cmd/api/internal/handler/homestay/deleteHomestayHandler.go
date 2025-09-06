package homestay

import (
	"golodge/common/result"
	"net/http"

	"golodge/app/travel/cmd/api/internal/logic/homestay"
	"golodge/app/travel/cmd/api/internal/svc"
	"golodge/app/travel/cmd/api/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func DeleteHomestayHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.DeleteHomestayReq
		if err := httpx.Parse(r, &req); err != nil {
			result.ParamErrorResult(r, w, err)
			return
		}

		l := homestay.NewDeleteHomestayLogic(r.Context(), svcCtx)
		resp, err := l.DeleteHomestay(&req)
		result.HttpResult(r, w, resp, err)
	}
}
