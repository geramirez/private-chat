package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

type MessageEvent struct {
	UserName  string
	Text      string
	TimeStamp int64
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

	chatService := NewChatService()
	webSocketServiceHub := NewWebSocketServiceHub()
	go webSocketServiceHub.Start()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, string(index))
	})

	http.HandleFunc("/chat", func(w http.ResponseWriter, r *http.Request) {

		webSocketService := NewWebSocketService(w, r, nil)
		defer webSocketService.websocketChannel.Close()

		fmt.Println("Connecting chat")

		for msg := range webSocketService.MessageStream() {
			chatService.Publish(msg)
			fmt.Println("Message from client to queue", string(msg))
		}
	})

	http.HandleFunc("/listen", func(w http.ResponseWriter, r *http.Request) {

		webSocketService := NewWebSocketService(w, r, nil)
		defer webSocketService.websocketChannel.Close()

		webSocketServiceHub.register <- webSocketService
		fmt.Println("Connecting listen")

		for d := range chatService.MessageStream() {
			webSocketServiceHub.broadcastChannel <- d.Body
			fmt.Println("Message from queue to client", string(d.Body))
		}
	})

	http.ListenAndServe(":3000", nil)
}
