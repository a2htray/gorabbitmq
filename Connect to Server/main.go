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

	// This may be a bug
	// @link https://github.com/rabbitmq/amqp091-go/issues/209
	connState := conn.ConnectionState()
	fmt.Println("Server name: ", connState.ServerName)
	fmt.Println("Version: ", connState.Version)

	channel, err := conn.Channel()
	if err != nil {
		log.Fatal(err)
	}

	defer channel.Close()

	fmt.Println("Channel is closed: ", channel.IsClosed())
}
