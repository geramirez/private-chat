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
			fmt.Println("closing channel")
		case message := <-clh.broadcastChannel:
			fmt.Println("Message from queue to clients", string(message))
			for client := range clh.clients {
				err := client.SendMessage(message)
				fmt.Println("total clients", len(clh.clients))
				if err != nil {
					fmt.Println("error closing channel")
					delete(clh.clients, client)
					client.websocketChannel.Close()
					fmt.Println("closing channel")
				}
			}
		}
	}
}

func (clh *WebSocketServiceHub) Unregister(ws *WebSocketService) {
	clh.unregister <- ws
}

func (clh *WebSocketServiceHub) Register(ws *WebSocketService) {
	clh.register <- ws
}

func (clh *WebSocketServiceHub) Publish(message []byte) {
	clh.broadcastChannel <- message
}
