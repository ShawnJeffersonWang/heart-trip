package homestayBussiness

import (
	"net/http"

	"heart-trip/app/travel/cmd/api/internal/logic/homestayBussiness"
	"heart-trip/app/travel/cmd/api/internal/svc"
	"heart-trip/app/travel/cmd/api/internal/types"
	"heart-trip/common/result"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func HomestayBussinessDetailHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.HomestayBussinessDetailReq
		if err := httpx.Parse(r, &req); err != nil {
			result.ParamErrorResult(r, w, err)
			return
		}

		l := homestayBussiness.NewHomestayBussinessDetailLogic(r.Context(), ctx)
		resp, err := l.HomestayBussinessDetail(req)
		result.HttpResult(r, w, resp, err)
	}
}
