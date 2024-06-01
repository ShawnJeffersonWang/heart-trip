package user

import (
	"context"
	"fmt"
	"github.com/gorilla/websocket"
	"golodge/app/websocket/cmd/api/internal/logic/ws"
	"log"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"
	"golodge/app/websocket/cmd/api/internal/svc"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type ChatLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// chat
func NewChatLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ChatLogic {
	return &ChatLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

type ChatRequest struct {
	UserId int64 `json:"userId"`
}

func (l *ChatLogic) Chat(w http.ResponseWriter, r *http.Request) {
	// todo: add your logic here and delete this line
	userID := "1"
	//userID := r.Header.Get("User-ID")
	//if userID == "" {
	//	http.Error(w, "Unauthorized", http.StatusUnauthorized)
	//	return
	//}

	fmt.Println("  --userID: ", userID)
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Chat: ", err)
		return
	}
	client := &ws.Client{
		Hub:    l.svcCtx.Hub,
		Conn:   conn,
		Send:   make(chan []byte, ws.BufSize),
		UserId: userID,
	}
	client.Hub.Register <- client

	go client.WritePump()
	go client.ReadPump()
}
