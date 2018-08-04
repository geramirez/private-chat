package main

import "fmt"

type WebSocketServiceHub struct {
	clients map[*WebSocketService]bool

	broadcastChannel chan []byte

	register chan *WebSocketService

	unregister chan *WebSocketService
}

func NewWebSocketServiceHub() *WebSocketServiceHub {
	return &WebSocketServiceHub{
		clients:          make(map[*WebSocketService]bool),
		broadcastChannel: make(chan []byte),
		register:         make(chan *WebSocketService),
		unregister:       make(chan *WebSocketService),
	}
}

func (clh *WebSocketServiceHub) Start() {
	for {
		select {
		case client := <-clh.register:
			clh.clients[client] = true
		case client := <-clh.unregister:
			delete(clh.clients, client)
			client.websocketChannel.Close()
		case message := <-clh.broadcastChannel:
			fmt.Println("message to be broadcast", string(message))
			for client := range clh.clients {
				client.SendMessage(message)
			}
		}
	}
}
