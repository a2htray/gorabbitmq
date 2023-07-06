package main

import (
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
)

func main() {
	conn, err := amqp.Dial("amqp://goadmin:123456@localhost:5672/goapp-vhost")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	channel, err := conn.Channel()
	if err != nil {
		log.Fatal(err)
	}
	defer channel.Close()

	// channel.QueueDeclare creates a queue if it doesn't exists, or ensures
	// that an existing queue matches the same properties.

	// a default binding to empty exchange "" which type is "direct" will be
	// created
	queue, err := channel.QueueDeclare(
		"goapp.test",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(queue)
}
