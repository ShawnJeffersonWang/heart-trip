package homestayComment

import (
	"golodge/common/result"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"golodge/app/travel/cmd/api/internal/logic/homestayComment"
	"golodge/app/travel/cmd/api/internal/svc"
	"golodge/app/travel/cmd/api/internal/types"
)

func AddCommentHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.AddCommentReq
		if err := httpx.Parse(r, &req); err != nil {
			result.ParamErrorResult(r, w, err)
			return
		}

		l := homestayComment.NewAddCommentLogic(r.Context(), svcCtx)
		resp, err := l.AddComment(&req)
		result.HttpResult(r, w, resp, err)
	}
}
