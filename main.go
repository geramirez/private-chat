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
		defer func() {
			fmt.Println("closing /chat")
			webSocketService.websocketChannel.Close()
		}()
		fmt.Println("Connecting chat")

		for wsMessage := range webSocketService.MessageStream() {
			if wsMessage.err != nil {
				break
			}
			chatService.Publish(wsMessage.message)
		}
	})

	http.HandleFunc("/listen", func(w http.ResponseWriter, r *http.Request) {

		webSocketService := NewWebSocketService(w, r, nil)
		defer func() {
			fmt.Println("closing /listen")
			webSocketServiceHub.Unregister(webSocketService)
		}()

		webSocketServiceHub.Register(webSocketService)
		fmt.Println("Connecting listen")

		for d := range chatService.MessageStream() {
			webSocketServiceHub.Publish(d.Body)
		}
	})

	http.ListenAndServe(":3000", nil)
}
