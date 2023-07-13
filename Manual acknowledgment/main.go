package main

import (
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
	FailedOnError(err, "connect to RabbitMQ server failed")
	defer conn.Close()

	channel, _ := conn.Channel()

	channel.Ack(1, true)
}
