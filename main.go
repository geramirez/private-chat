package main

import (
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

type ChatService struct {
	queue   amqp.Queue
	channel *amqp.Channel
}

func NewChatService() *ChatService {
	messageQueue, err := amqp.Dial("amqp://guest:guest@0.0.0.0:5672/")
	failOnError(err, "Failed to connect to RabbitMQ --")

	queueChannel, err := messageQueue.Channel()
	failOnError(err, "Failed to open a channel")

	queue, err := queueChannel.QueueDeclare(
		"chat", // name
		false,  // durable
		false,  // delete when unused
		false,  // exclusive
		false,  // no-wait
		nil,    // arguments
	)
	failOnError(err, "Failed to declare a queue")

	fmt.Println("New Chat Service...")

	return &ChatService{
		queue:   queue,
		channel: queueChannel,
	}
}

func (c *ChatService) Publish(message []byte) {
	err := c.channel.Publish(
		"",           // exchange
		c.queue.Name, // routing key
		false,        // mandatory
		false,        // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        message,
		})
	failOnError(err, "Failed to publish a message")
}

func (c *ChatService) GetNextMessage() <-chan amqp.Delivery {
	msgs, err := c.channel.Consume(
		c.queue.Name, // queue
		"",           // consumer
		true,         // auto-ack
		false,        // exclusive
		false,        // no-local
		false,        // no-wait
		nil,          // args
	)
	failOnError(err, "Failed to register a consumer")
	return msgs
}

type CabelService struct {
	websocketChannel *websocket.Conn
}

func NewCabelService(w http.ResponseWriter, r *http.Request, responseHeader http.Header) *CabelService {
	websocketChannel, err := buffer.Upgrade(w, r, nil)
	fmt.Println("New Cabel Service...")
	failOnError(err, "NewCabelService fail")
	return &CabelService{
		websocketChannel: websocketChannel,
	}
}

func (c *CabelService) ReadNextMessage() []byte {
	_, msg, err := c.websocketChannel.ReadMessage()
	failOnError(err, "WebsocketChannel Read Fail")
	return msg
}

func (c *CabelService) SendMessage(msg []byte) {
	err := c.websocketChannel.WriteMessage(websocket.TextMessage, msg)
	failOnError(err, "WebsocketChannel Send Fail")

}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
		panic(fmt.Sprintf("%s: %s", msg, err))
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

	chatService := NewChatService()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, string(index))
	})

	http.HandleFunc("/chat", func(w http.ResponseWriter, r *http.Request) {

		cabelService := NewCabelService(w, r, nil)

		for {
			msg := cabelService.ReadNextMessage()
			chatService.Publish(msg)
			fmt.Println("Message from client to queue", string(msg))
		}
	})

	http.HandleFunc("/listen", func(w http.ResponseWriter, r *http.Request) {

		cabelService := NewCabelService(w, r, nil)

		for d := range chatService.GetNextMessage() {
			cabelService.SendMessage(d.Body)
			fmt.Println("Message from queue to client", string(d.Body))
		}

	})

	http.ListenAndServe(":3000", nil)
}
