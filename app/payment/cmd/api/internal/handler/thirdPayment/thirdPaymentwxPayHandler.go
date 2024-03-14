package thirdPayment

import (
	"net/http"

	"homestay/app/payment/cmd/api/internal/logic/thirdPayment"
	"homestay/app/payment/cmd/api/internal/svc"
	"homestay/app/payment/cmd/api/internal/types"
	"homestay/common/result"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func ThirdPaymentwxPayHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ThirdPaymentWxPayReq
		if err := httpx.Parse(r, &req); err != nil {
			result.ParamErrorResult(r, w, err)
			return
		}

		l := thirdPayment.NewThirdPaymentwxPayLogic(r.Context(), ctx)
		resp, err := l.ThirdPaymentwxPay(req)
		result.HttpResult(r, w, resp, err)
	}
}
