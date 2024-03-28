package thirdPayment

import (
	"golodge/common/result"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"golodge/app/payment/cmd/api/internal/logic/thirdPayment"
	"golodge/app/payment/cmd/api/internal/svc"
	"golodge/app/payment/cmd/api/internal/types"
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
