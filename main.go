package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/gorilla/websocket"
)

type MessageEvent struct {
	UserName  string
	Text      string
	TimeStamp int64
}

var buffer = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func main() {
	indexFile, err := os.Open("index.html")
	if err != nil {
		fmt.Println(err)
	}
	index, err := ioutil.ReadAll(indexFile)
	if err != nil {
		fmt.Println(err)
	}

	chatClient := NewChatClient()
	webSocketServiceHub := NewWebSocketServiceHub()
	go webSocketServiceHub.Start()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, string(index))
	})

	http.HandleFunc("/chat", func(w http.ResponseWriter, r *http.Request) {

		webSocketService := NewWebSocketService(w, r, nil)
		defer webSocketService.websocketChannel.Close()

		fmt.Println("Connecting chat")

		for {
			msg := webSocketService.ReadNextMessage()
			chatClient.Publish(msg)
			fmt.Println("Message from client to queue", string(msg))
		}
	})

	http.HandleFunc("/listen", func(w http.ResponseWriter, r *http.Request) {

		webSocketService := NewWebSocketService(w, r, nil)
		defer webSocketService.websocketChannel.Close()

		webSocketServiceHub.register <- webSocketService
		fmt.Println("Connecting listen")

		for d := range chatClient.GetNextMessage() {
			webSocketServiceHub.broadcastChannel <- d.Body
			fmt.Println("Message from queue to client", string(d.Body))
		}

	})

	http.ListenAndServe(":3000", nil)
}
