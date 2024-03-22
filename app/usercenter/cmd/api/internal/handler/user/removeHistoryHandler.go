package user

import (
	"golodge/common/result"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"golodge/app/usercenter/cmd/api/internal/logic/user"
	"golodge/app/usercenter/cmd/api/internal/svc"
	"golodge/app/usercenter/cmd/api/internal/types"
)

func RemoveHistoryHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.RemoveHistoryReq
		if err := httpx.Parse(r, &req); err != nil {
			result.ParamErrorResult(r, w, err)
			return
		}

		l := user.NewRemoveHistoryLogic(r.Context(), svcCtx)
		resp, err := l.RemoveHistory(&req)
		result.HttpResult(r, w, resp, err)
	}
}
