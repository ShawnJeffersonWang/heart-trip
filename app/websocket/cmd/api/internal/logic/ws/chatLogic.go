package ws

import (
	"context"
	"github.com/gorilla/websocket"
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
	fromUserId := r.Header.Get("From-User-ID")
	toUserId := r.Header.Get("To-User-ID")

	if fromUserId == "" || toUserId == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	client := &Client{
		Hub:    l.svcCtx.Hub,
		Conn:   conn,
		Send:   make(chan []byte, BufSize),
		UserId: fromUserId,
	}
	client.Hub.Register <- client

	go client.WritePump()
	go client.ReadPump()
}
