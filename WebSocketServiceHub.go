package main

import "fmt"

type IWebSocketServiceHub interface {
	Start()
	Unregister(*WebSocketService)
	Register(*WebSocketService)
	Publish([]byte)
}

type WebSocketServiceHub struct {
	clients map[*WebSocketService]bool

	broadcastChannel chan []byte

	register chan *WebSocketService
}

func NewWebSocketServiceHub() *WebSocketServiceHub {
	return &WebSocketServiceHub{
		clients:          make(map[*WebSocketService]bool),
		broadcastChannel: make(chan []byte),
		register:         make(chan *WebSocketService),
	}
}

func (clh *WebSocketServiceHub) Start() {
	for {
		select {
		case client := <-clh.register:
			clh.clients[client] = true
		case message := <-clh.broadcastChannel:
			fmt.Println("Message from queue to" + string(len(clh.clients)) + "clients")
			for client := range clh.clients {
				err := client.SendMessage(message)
				if err != nil {
					delete(clh.clients, client)
					client.Close()
				}
			}
		}
	}
}

func (clh *WebSocketServiceHub) Unregister(client *WebSocketService) {
	delete(clh.clients, client)
	client.Close()
}

func (clh *WebSocketServiceHub) Register(ws *WebSocketService) {
	clh.register <- ws
}

func (clh *WebSocketServiceHub) Publish(message []byte) {
	clh.broadcastChannel <- message
}
