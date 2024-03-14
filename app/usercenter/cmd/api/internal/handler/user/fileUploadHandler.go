package user

import (
	"homestay/common/result"
	"net/http"

	"homestay/app/usercenter/cmd/api/internal/logic/user"
	"homestay/app/usercenter/cmd/api/internal/svc"
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
