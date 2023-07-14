package main

import (
	"fmt"
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

	channel, _ := conn.Channel()
	for i := 0; i < 10; i++ {
		// Those queues will not remain after RabbitMQ restart
		channel.QueueDeclare(
			fmt.Sprintf("short-lived.queue.%d", i),
			false,
			false,
			false,
			false,
			nil,
		)
	}
}
