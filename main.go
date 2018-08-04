package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/websocket"
	"github.com/streadway/amqp"
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

type msg struct {
	Num int
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
	messageQueue, err := amqp.Dial("amqp://guest:guest@0.0.0.0:5672/")
	failOnError(err, "Failed to connect to RabbitMQ --")
	defer messageQueue.Close()

	queueChannel, err := messageQueue.Channel()
	failOnError(err, "Failed to open a channel")
	defer queueChannel.Close()

	q, err := queueChannel.QueueDeclare(
		"chat", // name
		false,  // durable
		false,  // delete when unused
		false,  // exclusive
		false,  // no-wait
		nil,    // arguments
	)
	failOnError(err, "Failed to declare a queue")

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, string(index))
	})

	http.HandleFunc("/chat", func(w http.ResponseWriter, r *http.Request) {

		websocketChannel, err := buffer.Upgrade(w, r, nil)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("Listening for messages...")

		for {
			_, msg, err := websocketChannel.ReadMessage()
			if err != nil {
				fmt.Println(err)
				return
			}
			err = queueChannel.Publish(
				"",     // exchange
				q.Name, // routing key
				false,  // mandatory
				false,  // immediate
				amqp.Publishing{
					ContentType: "text/plain",
					Body:        msg,
				})
			failOnError(err, "Failed to publish a message")
			fmt.Println("Message read")
		}
	})

	http.HandleFunc("/listen", func(w http.ResponseWriter, r *http.Request) {

		websocketChannel, err := buffer.Upgrade(w, r, nil)
		if err != nil {
			fmt.Println(err)
			return
		}

		msgs, err := queueChannel.Consume(
			q.Name, // queue
			"",     // consumer
			true,   // auto-ack
			false,  // exclusive
			false,  // no-local
			false,  // no-wait
			nil,    // args
		)
		failOnError(err, "Failed to register a consumer")

		fmt.Println("Client setup")

		for d := range msgs {

			message := MessageEvent{}
			json.Unmarshal(d.Body, &message)
			fmt.Println(message)

			err = websocketChannel.WriteMessage(websocket.TextMessage, d.Body)
			if err != nil {
				fmt.Println(err)
				break
			}
		}

	})

	http.ListenAndServe(":3000", nil)
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
		panic(fmt.Sprintf("%s: %s", msg, err))
	}
}
