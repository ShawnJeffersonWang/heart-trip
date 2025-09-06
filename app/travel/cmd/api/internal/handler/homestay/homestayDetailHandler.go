package homestay

import (
	"heart-trip/app/travel/cmd/api/internal/logic/homestay"
	"net/http"

	"heart-trip/app/travel/cmd/api/internal/svc"
	"heart-trip/app/travel/cmd/api/internal/types"
	"heart-trip/common/result"

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
