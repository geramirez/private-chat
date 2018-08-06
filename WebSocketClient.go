package main

import (
	"net/http"

	"github.com/gorilla/websocket"
)

var buffer = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type IWebSocketClient interface {
	MessageStream() <-chan WebSocketMessage
	SendMessage([]byte) error
	Close()
}

type WebSocketClient struct {
	websocketChannel *websocket.Conn
}

type WebSocketMessage struct {
	message []byte
	err     error
}

func NewWebSocketClient(w http.ResponseWriter, r *http.Request, responseHeader http.Header) *WebSocketClient {
	websocketChannel, err := buffer.Upgrade(w, r, nil)
	failOnError(err, "NewWebSocketClient fail")
	return &WebSocketClient{
		websocketChannel: websocketChannel,
	}
}

func (c *WebSocketClient) MessageStream() <-chan WebSocketMessage {

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

func (c *WebSocketClient) readNextMessage() ([]byte, error) {
	_, msg, err := c.websocketChannel.ReadMessage()
	if err != nil {
		return nil, err
	}
	return msg, nil
}

func (c *WebSocketClient) SendMessage(msg []byte) error {
	err := c.websocketChannel.WriteMessage(websocket.TextMessage, msg)
	if err != nil {
		return err
	}
	return nil
}

func (c *WebSocketClient) Close() {
	c.websocketChannel.Close()
}
