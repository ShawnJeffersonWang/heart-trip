package homestay

import (
	"golodge/common/result"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"golodge/app/travel/cmd/api/internal/logic/homestay"
	"golodge/app/travel/cmd/api/internal/svc"
	"golodge/app/travel/cmd/api/internal/types"
)

func AddHomestayHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.AddHomestayReq
		if err := httpx.Parse(r, &req); err != nil {
			result.ParamErrorResult(r, w, err)
			return
		}

		l := homestay.NewAddHomestayLogic(r.Context(), svcCtx)
		resp, err := l.AddHomestay(&req)
		result.HttpResult(r, w, resp, err)
	}
}
