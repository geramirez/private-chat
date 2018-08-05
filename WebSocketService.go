package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

var buffer = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type WebSocketService struct {
	websocketChannel *websocket.Conn
}

type WebSocketMessage struct {
	message []byte
	err     error
}

func NewWebSocketService(w http.ResponseWriter, r *http.Request, responseHeader http.Header) *WebSocketService {
	websocketChannel, err := buffer.Upgrade(w, r, nil)
	fmt.Println("New Cabel Service...")
	failOnError(err, "NewWebSocketService fail")
	return &WebSocketService{
		websocketChannel: websocketChannel,
	}
}

func (c *WebSocketService) MessageStream() <-chan WebSocketMessage {

	messageChannel := make(chan WebSocketMessage)
	go func() {
		for {
			message, err := c.readNextMessage()

			messageChannel <- WebSocketMessage{
				message: message,
				err:     err,
			}
		}
	}()
	return messageChannel
}

func (c *WebSocketService) readNextMessage() ([]byte, error) {
	msgType, msg, err := c.websocketChannel.ReadMessage()
	if err != nil {
		fmt.Println(string(msgType), string(msg), err)
		return nil, err
	}
	return msg, nil
}

func (c *WebSocketService) SendMessage(msg []byte) error {
	err := c.websocketChannel.WriteMessage(websocket.TextMessage, msg)
	if err != nil {
		return err
	}
	return nil
}

func (c *WebSocketService) Close() {
	fmt.Println("closing channel")
	c.websocketChannel.Close()
}
