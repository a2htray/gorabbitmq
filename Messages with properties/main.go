package main

import (
	"context"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"time"
)

func FailedOnError(err error, message string) {
	if err != nil {
		log.Fatal(message, "Error: ", err.Error())
	}
}

func main() {
	conn, err := amqp.Dial("amqp://goadmin:123456@localhost:5672/goapp-vhost")
	FailedOnError(err, "could not connect to RabbitMQ server")
	defer conn.Close()

	channel, err := conn.Channel()
	FailedOnError(err, "could not create a channel from a connection")
	defer channel.Close()

	queue, err := channel.QueueDeclare(
		"goapp.queue.test",
		false,
		false,
		false,
		false,
		nil,
	)
	FailedOnError(err, "declare queue failed")

	err = channel.ExchangeDeclare(
		"goapp.exchange.direct",
		amqp.ExchangeDirect,
		true,
		false,
		false,
		true,
		nil,
	)
	FailedOnError(err, "declare exchange failed")

	routingKey := "routing.key"
	err = channel.QueueBind(queue.Name, routingKey, "goapp.exchange.direct", false, nil)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	for i := 0; i < 10; i++ {
		message := amqp.Publishing{
			Headers:         nil,
			ContentType:     "plain/text",
			ContentEncoding: "UTF-8",
			DeliveryMode:    2, // Transient (0 or 1) or Persistent (2)
			Priority:        uint8(i),
			MessageId:       fmt.Sprintf("message%d", i),
			Timestamp:       time.Now(),
			Type:            "gotest",
			UserId:          "goadmin",
			AppId:           "Messages with properties",
			Body:            []byte(fmt.Sprintf("hello user%d", i)),
		}
		err = channel.PublishWithContext(
			ctx,
			"goapp.exchange.direct",
			routingKey,
			false,
			false,
			message,
		)

		FailedOnError(err, "publish message failed")
	}
}
