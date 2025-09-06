package homestayBussiness

import (
	"net/http"

	"heart-trip/app/travel/cmd/api/internal/logic/homestayBussiness"
	"heart-trip/app/travel/cmd/api/internal/svc"
	"heart-trip/app/travel/cmd/api/internal/types"
	"heart-trip/common/result"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func GoodBossHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GoodBossReq
		if err := httpx.Parse(r, &req); err != nil {
			result.ParamErrorResult(r, w, err)
			return
		}

		l := homestayBussiness.NewGoodBossLogic(r.Context(), ctx)
		resp, err := l.GoodBoss(req)
		result.HttpResult(r, w, resp, err)
	}
}
