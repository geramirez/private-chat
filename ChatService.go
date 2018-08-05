package main

import (
	"fmt"

	"github.com/streadway/amqp"
)

type ChatService struct {
	queue   amqp.Queue
	channel *amqp.Channel
}

type IChatService interface {
	Publish([]byte)
	MessageStream() <-chan []byte
}

func NewChatService(queueUrl string) *ChatService {
	messageQueue, err := amqp.Dial(queueUrl)
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
	fmt.Println("Message from client to queue", string(message))
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

func (c *ChatService) MessageStream() <-chan []byte {
	queueMessageChannel, err := c.channel.Consume(
		c.queue.Name, // queue
		"",           // consumer
		true,         // auto-ack
		false,        // exclusive
		false,        // no-local
		false,        // no-wait
		nil,          // args
	)
	failOnError(err, "Failed to register a consumer")

	messageChannel := make(chan []byte)
	go func() {
		for payload := range queueMessageChannel {
			messageChannel <- payload.Body
		}
	}()
	return messageChannel
}
