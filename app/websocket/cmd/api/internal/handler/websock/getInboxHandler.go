package websock

import (
	"heart-trip/common/result"
	"net/http"

	"heart-trip/app/websocket/cmd/api/internal/logic/websock"
	"heart-trip/app/websocket/cmd/api/internal/svc"
	"heart-trip/app/websocket/cmd/api/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

// get inbox messages
func GetInboxHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetInboxReq
		if err := httpx.Parse(r, &req); err != nil {
			result.ParamErrorResult(r, w, err)
			return
		}

		l := websock.NewGetInboxLogic(r.Context(), svcCtx)
		resp, err := l.GetInbox(&req)
		result.HttpResult(r, w, resp, err)
	}
}
