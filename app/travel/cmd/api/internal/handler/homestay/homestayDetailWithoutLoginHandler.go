package homestay

import (
	"github.com/zeromicro/go-zero/rest/httpx"
	"golodge/app/travel/cmd/api/internal/logic/homestay"
	"golodge/app/travel/cmd/api/internal/svc"
	"golodge/app/travel/cmd/api/internal/types"
	"golodge/common/result"
	"net/http"
)

//var domain = flag.String("domain", "http://localhost:8888", "the domain to request")

func HomestayDetailWithoutLoginHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.HomestayDetailReq
		//fmt.Println("Authorization != \"\"")
		if err := httpx.Parse(r, &req); err != nil {
			result.ParamErrorResult(r, w, err)
			return
		}
		//if len(r.Header.Get("Authorization")) > 0 {
		//	l := homestay.NewHomestayDetailLogic(r.Context(), svcCtx)
		//	resp, err := l.HomestayDetail(req)
		//	result.HttpResult(r, w, resp, err)
		//	//response, err := httpc.Do(r.Context(), http.MethodPost, *domain+"/travel/v1/homestay/homestayDetail", req)
		//	//if err != nil {
		//	//	fmt.Println("httpc.Do err", err)
		//	//	return
		//	//}
		//	//io.Copy(os.Stdout, response.Body)
		//	return
		//}
		l := homestay.NewHomestayDetailWithoutLoginLogic(r.Context(), svcCtx)
		resp, err := l.HomestayDetailWithoutLogin(&req)
		result.HttpResult(r, w, resp, err)
	}
}
