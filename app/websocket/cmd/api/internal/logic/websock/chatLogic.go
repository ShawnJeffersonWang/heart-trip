package websock

import (
	"context"
	"github.com/gorilla/websocket"
	"github.com/zeromicro/go-zero/core/logx"
	"golodge/app/websocket/cmd/api/internal/svc"
	"golodge/common/ctxdata"
	"net/http"
	"strconv"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
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

func (l *ChatLogic) Chat(w http.ResponseWriter, r *http.Request, hub *svc.Hub) {
	// todo: add your logic here and delete this line
	// bug: 不能fromUserId := string(ctxdata.GetUidFromCtx(l.ctx))
	fromUserId := strconv.FormatInt(ctxdata.GetUidFromCtx(l.ctx), 10)
	//toUserId := r.Header.Get("To-User-ID")
	toUserId := r.URL.Query().Get("receiver")

	if fromUserId == "" || toUserId == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	client := &svc.Client{
		Hub:    hub,
		Conn:   conn,
		Send:   make(chan []byte, svc.BufSize),
		UserId: fromUserId,
	}
	client.Hub.Register <- client

	go client.WritePump(fromUserId, toUserId)
	go client.ReadPump(fromUserId, toUserId, l.svcCtx)
}
