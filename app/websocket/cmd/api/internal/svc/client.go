package svc

import (
	"bytes"
	"context"
	"fmt"
	"heart-trip/app/websocket/cmd/api/internal/types"
	"heart-trip/app/websocket/model"
	"log"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 5 * time.Minute
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 512
	BufSize        = 256
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

type Client struct {
	Hub    *Hub
	Conn   *websocket.Conn
	Send   chan []byte
	UserId string
}

func (c *Client) ReadPump(fromUserId, toUserId string, svcCtx *ServiceContext) {
	defer func() {
		c.Hub.Unregister <- c
		c.Conn.Close()
		log.Println("ReadPump.defer func()")
	}()
	c.Conn.SetReadLimit(maxMessageSize)
	c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	c.Conn.SetPongHandler(func(string) error {
		_ = c.Conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})
	for {
		//err := c.Conn.ReadJSON(&msg)
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		log.Printf("ReadPump fromUserId: %s, toUserId: %s\n", fromUserId, toUserId)
		message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
		msg := types.Message{
			ToUserId:   toUserId,
			FromUserId: fromUserId,
			Content:    string(message),
			Type:       "1",
		}
		c.Hub.Broadcast <- msg
		fromId, _ := strconv.Atoi(fromUserId)
		toId, _ := strconv.Atoi(toUserId)
		insertMsg := model.Message{
			FromUserId: int64(fromId),
			ToUserId:   int64(toId),
			Content:    string(message),
		}
		_, err = svcCtx.MessageModel.Insert(context.Background(), &insertMsg)
		if err != nil {
			fmt.Println("插入失败")
			return
		}
	}
}

func (c *Client) WritePump(fromUserId, toUserId string) {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
		log.Println("WritePump.defer func()")
	}()

	for {
		select {
		case msg, ok := <-c.Send:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			log.Println("WritePump.msg: ", string(msg))
			c.Conn.WriteMessage(websocket.TextMessage, msg)
			//w, err := c.Conn.NextWriter(websocket.TextMessage)
			//if err != nil {
			//	return
			//}
			//w.Write(msg)
			//
			//n := len(c.Send)
			//for i := 0; i < n; i++ {
			//	w.Write(newline)
			//	w.Write(<-c.Send)
			//}
			//if err := w.Close(); err != nil {
			//	return
			//}
		case <-ticker.C:
			log.Println("WritePump.ticker.C", ticker.C)
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
