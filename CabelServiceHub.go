package main

import "fmt"

type CabelServiceHub struct {
	clients map[*CabelService]bool

	broadcastChannel chan []byte

	register chan *CabelService

	unregister chan *CabelService
}

func NewCabelServiceHub() *CabelServiceHub {
	return &CabelServiceHub{
		clients:          make(map[*CabelService]bool),
		broadcastChannel: make(chan []byte),
		register:         make(chan *CabelService),
		unregister:       make(chan *CabelService),
	}
}

func (clh *CabelServiceHub) Start() {
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
