package main

import (
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

	err = channel.ExchangeDeclare(
		"goapp.exchange.direct",
		amqp.ExchangeDirect,
		false,
		false,
		false,
		true,
		nil,
	)
	if err != nil {
		log.Fatal(err)
	}
}
