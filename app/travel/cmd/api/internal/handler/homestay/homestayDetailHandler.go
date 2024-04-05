package homestay

import (
	"golodge/app/travel/cmd/api/internal/logic/homestay"
	"net/http"

	"golodge/app/travel/cmd/api/internal/svc"
	"golodge/app/travel/cmd/api/internal/types"
	"golodge/common/result"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func HomestayDetailHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.HomestayDetailReq
		//getResult := r.Header.Get("Authorization")
		//fmt.Println("Authorization: ", getResult)
		if err := httpx.Parse(r, &req); err != nil {
			result.ParamErrorResult(r, w, err)
			return
		}

		l := homestay.NewHomestayDetailLogic(r.Context(), ctx)
		resp, err := l.HomestayDetail(req)
		result.HttpResult(r, w, resp, err)
	}
}
