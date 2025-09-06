package svc

import (
	"encoding/json"
	"heart-trip/app/websocket/cmd/api/internal/types"
	"log"
)

type Hub struct {
	Clients    map[string]*Client
	Broadcast  chan types.Message
	Register   chan *Client
	Unregister chan *Client
}

func NewHub() *Hub {
	return &Hub{
		Clients:    make(map[string]*Client),
		Broadcast:  make(chan types.Message),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.Clients[client.UserId] = client
			log.Println("Run.Register: ", client)
			//h.sendOnlineUserList()
		case client := <-h.Unregister:
			if _, ok := h.Clients[client.UserId]; ok {
				delete(h.Clients, client.UserId)
				close(client.Send)
				log.Println("Run.Unregister: ", client.UserId)
			}
			//h.sendOnlineUserList()
		case message := <-h.Broadcast:
			if receiver, ok := h.Clients[message.ToUserId]; ok {
				log.Println("Run.Broadcast: ", "receiver.UserId: ", receiver.UserId)
				receiver.Send <- []byte(message.Content)
			}
		}
	}
}

func (h *Hub) sendOnlineUserList() {
	var connectList []string
	for userId := range h.Clients {
		connectList = append(connectList, userId)
		bytes, _ := json.Marshal(connectList)
		h.SendBroadCastMessage(nil, bytes, false)
	}
}

//func (h *Hub) SendMessage(message []byte) {
//	if msg.ToUserId != "" {
//		if client, ok := h.Clients[msg.ToUserId]; ok {
//			client.Send <- message
//		}
//	} else {
//		h.SendBroadCastMessage(msg, message, true)
//	}
//}

func (h *Hub) SendBroadCastMessage(msg *types.Message, message []byte, skipSelf bool) {
	for _, client := range h.Clients {
		if skipSelf && client.UserId == msg.FromUserId {
			continue
		}
		select {
		case client.Send <- message:
		default:
			close(client.Send)
			delete(h.Clients, client.UserId)
		}
	}
}
