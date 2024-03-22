package homestayOrder

import (
	"net/http"

	"golodge/app/order/cmd/api/internal/logic/homestayOrder"
	"golodge/app/order/cmd/api/internal/svc"
	"golodge/app/order/cmd/api/internal/types"
	"golodge/common/result"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func CreateHomestayOrderHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.CreateHomestayOrderReq
		if err := httpx.Parse(r, &req); err != nil {
			result.ParamErrorResult(r, w, err)
			return
		}

		l := homestayOrder.NewCreateHomestayOrderLogic(r.Context(), ctx)
		resp, err := l.CreateHomestayOrder(req)
		result.HttpResult(r, w, resp, err)
	}
}
