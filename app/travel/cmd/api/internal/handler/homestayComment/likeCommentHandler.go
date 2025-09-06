package homestayComment

import (
	"heart-trip/common/result"
	"net/http"

	"heart-trip/app/travel/cmd/api/internal/logic/homestayComment"
	"heart-trip/app/travel/cmd/api/internal/svc"
	"heart-trip/app/travel/cmd/api/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func LikeCommentHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.LikeCommentReq
		if err := httpx.Parse(r, &req); err != nil {
			result.ParamErrorResult(r, w, err)
			return
		}

		l := homestayComment.NewLikeCommentLogic(r.Context(), svcCtx)
		resp, err := l.LikeComment(&req)
		result.HttpResult(r, w, resp, err)
	}
}
