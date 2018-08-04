package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

type WebSocketService struct {
	websocketChannel *websocket.Conn
}

func NewWebSocketService(w http.ResponseWriter, r *http.Request, responseHeader http.Header) *WebSocketService {
	websocketChannel, err := buffer.Upgrade(w, r, nil)
	fmt.Println("New Cabel Service...")
	failOnError(err, "NewWebSocketService fail")
	return &WebSocketService{
		websocketChannel: websocketChannel,
	}
}

func (c *WebSocketService) ReadNextMessage() []byte {
	msgType, msg, err := c.websocketChannel.ReadMessage()
	fmt.Println(string(msgType), string(msg), err)
	failOnError(err, "WebsocketChannel Read Fail")
	return msg
}

func (c *WebSocketService) SendMessage(msg []byte) {
	err := c.websocketChannel.WriteMessage(websocket.TextMessage, msg)
	failOnError(err, "WebsocketChannel Send Fail")
}
