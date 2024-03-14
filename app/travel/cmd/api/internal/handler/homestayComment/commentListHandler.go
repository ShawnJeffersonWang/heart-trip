package homestayComment

import (
	"net/http"

	"homestay/app/travel/cmd/api/internal/logic/homestayComment"
	"homestay/app/travel/cmd/api/internal/svc"
	"homestay/app/travel/cmd/api/internal/types"
	"homestay/common/result"

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
