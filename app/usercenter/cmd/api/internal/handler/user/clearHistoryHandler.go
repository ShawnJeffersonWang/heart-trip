package user

import (
	"golodge/common/result"
	"net/http"

	"golodge/app/usercenter/cmd/api/internal/logic/user"
	"golodge/app/usercenter/cmd/api/internal/svc"
	"golodge/app/usercenter/cmd/api/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func ClearHistoryHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ClearHistoryReq
		if err := httpx.Parse(r, &req); err != nil {
			result.ParamErrorResult(r, w, err)
			return
		}

		l := user.NewClearHistoryLogic(r.Context(), svcCtx)
		resp, err := l.ClearHistory(&req)
		result.HttpResult(r, w, resp, err)
	}
}
