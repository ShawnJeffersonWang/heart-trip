package thirdPayment

import (
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpx"
	"golodge/app/payment/cmd/api/internal/types"
	"golodge/common/result"
	"net/http"

	"golodge/app/payment/cmd/api/internal/logic/thirdPayment"
	"golodge/app/payment/cmd/api/internal/svc"
)

func FakePayCallbackHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.FakePayCallbackReq
		if err := httpx.Parse(r, &req); err != nil {
			result.ParamErrorResult(r, w, err)
			return
		}

		l := thirdPayment.NewFakePayCallbackLogic(r.Context(), svcCtx)
		resp, err := l.FakePayCallback(&req)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
		} else {
			w.WriteHeader(http.StatusOK)
		}
		logx.Infof("ReturnCode: %s", resp.ReturnCode)
		//fmt.Fprint(w, resp.ReturnCode)
		result.HttpResult(r, w, resp, err)
	}
}
