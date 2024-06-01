package user

import (
	"net/http"

	"golodge/app/websocket/cmd/api/internal/logic/user"
	"golodge/app/websocket/cmd/api/internal/svc"
)

// chat
func ChatHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := user.NewChatLogic(r.Context(), svcCtx)
		l.Chat(w, r)
	}
}
