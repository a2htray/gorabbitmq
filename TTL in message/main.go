package main

import (
	"context"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
)

func FailedOnError(err error, message string) {
	if err != nil {
		log.Fatal(message, "Error: ", err.Error())
	}
}

func main() {
	conn, err := amqp.Dial("amqp://goadmin:123456@localhost:5672/goapp-vhost")
	FailedOnError(err, "failed to create connection to RabbitMQ server")
	defer conn.Close()

	channel, err := conn.Channel()
	FailedOnError(err, "failed to open a channel")

	queue, err := channel.QueueDeclare(
		"goapp.queue.ttl.10s",
		true,
		false,
		false,
		false,
		nil,
	)

	_ = channel.PublishWithContext(
		context.Background(),
		"",
		queue.Name,
		false,
		false,
		amqp.Publishing{
			Expiration:  "10000",
			ContentType: "text/plain",
			Body:        []byte("message with ttl 10000"),
		},
	)

	select {}
}
