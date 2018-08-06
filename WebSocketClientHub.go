package main

type IWebSocketClientHub interface {
	Start()
	Unregister(IWebSocketClient)
	Register(IWebSocketClient)
	Publish([]byte)
}

type WebSocketServiceHub struct {
	clients        map[IWebSocketClient]bool
	broadcast      chan []byte
	broadcastDone  chan bool
	register       chan IWebSocketClient
	registerDone   chan bool
	unregister     chan IWebSocketClient
	unregisterDone chan bool
}

func NewWebSocketClientHub() *WebSocketServiceHub {
	return &WebSocketServiceHub{
		clients:        make(map[IWebSocketClient]bool),
		broadcast:      make(chan []byte),
		broadcastDone:  make(chan bool),
		register:       make(chan IWebSocketClient),
		registerDone:   make(chan bool),
		unregister:     make(chan IWebSocketClient),
		unregisterDone: make(chan bool),
	}
}

func (h *WebSocketServiceHub) Start() {
	for {
		select {
		case client := <-h.register:
			h.addClient(client)
			h.registerDone <- true
		case client := <-h.unregister:
			h.removeClient(client)
			h.unregisterDone <- true
		case message := <-h.broadcast:
			h.broadcastToAll(message)
			h.broadcastDone <- true
		}
	}
}

func (h *WebSocketServiceHub) Unregister(client IWebSocketClient) {
	h.unregister <- client
	<-h.unregisterDone
}

func (h *WebSocketServiceHub) Register(client IWebSocketClient) {
	h.register <- client
	<-h.registerDone
}

func (h *WebSocketServiceHub) Publish(message []byte) {
	h.broadcast <- message
	<-h.broadcastDone
}

func (h *WebSocketServiceHub) removeClient(client IWebSocketClient) {
	delete(h.clients, client)
	client.Close()
}

func (h *WebSocketServiceHub) addClient(client IWebSocketClient) {
	h.clients[client] = true
}

func (h *WebSocketServiceHub) broadcastToAll(message []byte) {
	for client := range h.clients {
		err := client.SendMessage(message)
		if err != nil {
			h.removeClient(client)
		}
	}
}
