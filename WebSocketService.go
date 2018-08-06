package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

var buffer = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type IWebSocketService interface {
	MessageStream() <-chan WebSocketMessage
	SendMessage([]byte) error
	Close()
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

//////////////////////////////////
type FakeWebSocketService struct {
	SendMessageArgs [][]byte
	CloseCalled     bool
}

func NewFakeWebSocketService() *FakeWebSocketService {
	return &FakeWebSocketService{}
}

func (c *FakeWebSocketService) MessageStream() <-chan WebSocketMessage {
	return make(chan WebSocketMessage)
}

func (c *FakeWebSocketService) SendMessage(msg []byte) error {
	c.SendMessageArgs = append(c.SendMessageArgs, msg)
	return nil
}

func (c *FakeWebSocketService) Close() {
	c.CloseCalled = true
	fmt.Println("Closing client")
}

////
type FakeWebBrokenSocketService struct {
	CloseCalled bool
}

func NewFakeWebBrokenSocketService() *FakeWebBrokenSocketService {
	return &FakeWebBrokenSocketService{}
}

func (c *FakeWebBrokenSocketService) MessageStream() <-chan WebSocketMessage {
	return make(chan WebSocketMessage)
}

func (c *FakeWebBrokenSocketService) SendMessage(msg []byte) error {
	return errors.New("")
}

func (c *FakeWebBrokenSocketService) Close() {
	c.CloseCalled = true
	fmt.Println("Closing client")
}
