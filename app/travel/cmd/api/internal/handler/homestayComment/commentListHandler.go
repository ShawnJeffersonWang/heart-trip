package homestayComment

import (
	"net/http"

	"heart-trip/app/travel/cmd/api/internal/logic/homestayComment"
	"heart-trip/app/travel/cmd/api/internal/svc"
	"heart-trip/app/travel/cmd/api/internal/types"
	"heart-trip/common/result"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func CommentListHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.CommentListReq
		if err := httpx.Parse(r, &req); err != nil {
			result.ParamErrorResult(r, w, err)
			return
		}

		l := homestayComment.NewCommentListLogic(r.Context(), ctx)
		resp, err := l.CommentList(req)
		result.HttpResult(r, w, resp, err)
	}
}
