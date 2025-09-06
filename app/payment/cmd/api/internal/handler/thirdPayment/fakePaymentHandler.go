package thirdPayment

import (
	"heart-trip/common/result"
	"net/http"

	"heart-trip/app/payment/cmd/api/internal/logic/thirdPayment"
	"heart-trip/app/payment/cmd/api/internal/svc"
	"heart-trip/app/payment/cmd/api/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func FakePaymentHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.FakePaymentReq
		if err := httpx.Parse(r, &req); err != nil {
			result.ParamErrorResult(r, w, err)
			return
		}

		l := thirdPayment.NewFakePaymentLogic(r.Context(), svcCtx)
		resp, err := l.FakePayment(&req)
		result.HttpResult(r, w, resp, err)
	}
}
