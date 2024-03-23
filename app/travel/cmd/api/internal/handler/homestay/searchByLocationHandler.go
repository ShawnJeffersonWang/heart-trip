package homestay

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"golodge/app/travel/cmd/api/internal/logic/homestay"
	"golodge/app/travel/cmd/api/internal/svc"
	"golodge/app/travel/cmd/api/internal/types"
)

func SearchByLocationHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.SearchByLocationReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := homestay.NewSearchByLocationLogic(r.Context(), svcCtx)
		resp, err := l.SearchByLocation(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
