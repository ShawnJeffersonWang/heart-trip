package websock

import (
	"golodge/common/result"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"golodge/app/websocket/cmd/api/internal/logic/websock"
	"golodge/app/websocket/cmd/api/internal/svc"
	"golodge/app/websocket/cmd/api/internal/types"
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
