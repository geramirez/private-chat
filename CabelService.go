package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

type CabelService struct {
	websocketChannel *websocket.Conn
}

func NewCabelService(w http.ResponseWriter, r *http.Request, responseHeader http.Header) *CabelService {
	websocketChannel, err := buffer.Upgrade(w, r, nil)
	fmt.Println("New Cabel Service...")
	failOnError(err, "NewCabelService fail")
	return &CabelService{
		websocketChannel: websocketChannel,
	}
}

func (c *CabelService) ReadNextMessage() []byte {
	msgType, msg, err := c.websocketChannel.ReadMessage()
	fmt.Println(string(msgType), string(msg), err)
	failOnError(err, "WebsocketChannel Read Fail")
	return msg
}

func (c *CabelService) SendMessage(msg []byte) {
	err := c.websocketChannel.WriteMessage(websocket.TextMessage, msg)
	failOnError(err, "WebsocketChannel Send Fail")

}
