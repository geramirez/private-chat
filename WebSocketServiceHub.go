package main

type IWebSocketServiceHub interface {
	Start()
	Unregister(IWebSocketService)
	Register(IWebSocketService)
	Publish([]byte)
}

type WebSocketServiceHub struct {
	clients        map[IWebSocketService]bool
	broadcast      chan []byte
	broadcastDone  chan bool
	register       chan IWebSocketService
	registerDone   chan bool
	unregister     chan IWebSocketService
	unregisterDone chan bool
}

func NewWebSocketServiceHub() *WebSocketServiceHub {
	return &WebSocketServiceHub{
		clients:        make(map[IWebSocketService]bool),
		broadcast:      make(chan []byte),
		broadcastDone:  make(chan bool),
		register:       make(chan IWebSocketService),
		registerDone:   make(chan bool),
		unregister:     make(chan IWebSocketService),
		unregisterDone: make(chan bool),
	}
}

func (clh *WebSocketServiceHub) Start() {
	for {
		select {
		case client := <-clh.register:
			clh.addClient(client)
			clh.registerDone <- true
		case client := <-clh.unregister:
			clh.removeClient(client)
			clh.unregisterDone <- true
		case message := <-clh.broadcast:
			clh.broadcastToAll(message)
			clh.broadcastDone <- true
		}
	}
}

func (clh *WebSocketServiceHub) Unregister(ws IWebSocketService) {
	clh.unregister <- ws
	<-clh.unregisterDone
}

func (clh *WebSocketServiceHub) Register(ws IWebSocketService) {
	clh.register <- ws
	<-clh.registerDone
}

func (clh *WebSocketServiceHub) Publish(message []byte) {
	clh.broadcast <- message
	<-clh.broadcastDone
}

func (clh *WebSocketServiceHub) removeClient(ws IWebSocketService) {
	delete(clh.clients, ws)
	ws.Close()
}

func (clh *WebSocketServiceHub) addClient(client IWebSocketService) {
	clh.clients[client] = true
}

func (clh *WebSocketServiceHub) broadcastToAll(message []byte) {
	for client := range clh.clients {
		err := client.SendMessage(message)
		if err != nil {
			clh.removeClient(client)
		}
	}
}
