package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	cfenv "github.com/cloudfoundry-community/go-cfenv"
)

type MessageEvent struct {
	UserName  string
	Text      string
	TimeStamp int64
}

type Config struct {
	Port     int
	QueueUrl string
}

type RabbitMQ struct {
}

func NewConfig() *Config {
	appEnv, err := cfenv.Current()
	if err != nil {
		return &Config{
			Port:     3000,
			QueueUrl: "amqp://guest:guest@0.0.0.0:5672/",
		}
	}
	services, _ := appEnv.Services.WithTag("RabbitMQ")
	rabbitMQURI := services[0].Credentials["uri"]
	return &Config{
		Port:     appEnv.Port,
		QueueUrl: rabbitMQURI.(string),
	}
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

	config := NewConfig()
	chatService := NewChatService(config.QueueUrl)
	webSocketClientHub := NewWebSocketClientHub()
	go webSocketClientHub.Start()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, string(index))
	})

	http.HandleFunc("/chat", func(w http.ResponseWriter, r *http.Request) {

		webSocketClient := NewWebSocketClient(w, r, nil)
		defer func() {
			fmt.Println("closing /chat")
			webSocketClient.websocketChannel.Close()
		}()
		fmt.Println("Connecting chat")

		for wsMessage := range webSocketClient.MessageStream() {
			if wsMessage.err != nil {
				break
			}
			chatService.Publish(wsMessage.message)
		}
	})

	http.HandleFunc("/listen", func(w http.ResponseWriter, r *http.Request) {

		webSocketClient := NewWebSocketClient(w, r, nil)
		defer func() {
			fmt.Println("closing /listen")
			webSocketClientHub.Unregister(webSocketClient)
		}()

		webSocketClientHub.Register(webSocketClient)
		fmt.Println("Connecting listen")

		for message := range chatService.MessageStream() {
			webSocketClientHub.Publish(message)
		}
	})
	fmt.Println("RUNNING ONT PORT:", fmt.Sprintf(":%d", config.Port))
	http.ListenAndServe(fmt.Sprintf(":%d", config.Port), nil)
}
