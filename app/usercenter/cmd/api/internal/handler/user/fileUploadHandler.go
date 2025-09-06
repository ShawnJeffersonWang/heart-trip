package user

import (
	"heart-trip/common/result"
	"net/http"

	"heart-trip/app/usercenter/cmd/api/internal/logic/user"
	"heart-trip/app/usercenter/cmd/api/internal/svc"
)

func FileUploadHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := user.NewFileUploadLogic(r.Context(), svcCtx)
		resp, err := l.FileUpload(r)
		if err != nil {
			result.ParamErrorResult(r, w, err)
		} else {
			result.HttpResult(r, w, resp, err)
		}
	}
}
