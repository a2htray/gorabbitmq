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

	fmt.Println("Connect is closed: ", conn.IsClosed())

	channel, err := conn.Channel()
	if err != nil {
		log.Fatal(err)
	}

	defer channel.Close()

	fmt.Println("Channel is closed: ", channel.IsClosed())
}
