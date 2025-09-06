package websock

import (
	"heart-trip/app/websocket/cmd/api/internal/logic/websock"
	"heart-trip/app/websocket/cmd/api/internal/svc"
	"net/http"
)

// chat
func ChatHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := websock.NewChatLogic(r.Context(), svcCtx)
		l.Chat(w, r, svcCtx.Hub)
	}
}
