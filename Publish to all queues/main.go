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
	FailedOnError(err, "failed to create connection to RabbitMQ server")
	defer conn.Close()

	channel, _ := conn.Channel()
	_ = channel.ExchangeDeclare(
		"goapp.exchange.fanout",
		amqp.ExchangeFanout,
		true,
		false,
		false,
		false,
		nil,
	)

	queue1, err := channel.QueueDeclare(
		"goapp.queue.user1",
		true,
		false,
		false,
		false,
		nil,
	)

	queue2, err := channel.QueueDeclare(
		"goapp.queue.user2",
		true,
		false,
		false,
		false,
		nil,
	)

	channel.QueueBind(
		queue1.Name,
		"info",
		"goapp.exchange.fanout",
		false,
		nil,
	)
	channel.QueueBind(
		queue2.Name,
		"info",
		"goapp.exchange.fanout",
		false,
		nil,
	)

	go func() {
		ctx := context.Background()
		for i := 0; i < 10; i++ {
			channel.PublishWithContext(
				ctx,
				"goapp.exchange.fanout",
				"info",
				false,
				false,
				amqp.Publishing{
					ContentType:  "plain/text",
					DeliveryMode: amqp.Persistent,
					Body:         []byte(fmt.Sprintf("message id: %d, time: %v", i+1, time.Now())),
				},
			)
		}
	}()

	go func() {
		userChannel1, _ := conn.Channel()
		chanDelivery, _ := userChannel1.Consume(
			queue1.Name,
			"user1",
			false,
			true,
			false,
			false,
			nil,
		)
		for delivery := range chanDelivery {
			fmt.Printf("user1 receives message: <%s>\n", string(delivery.Body))
			_ = delivery.Ack(false)
		}
	}()

	go func() {
		userChannel2, _ := conn.Channel()
		chanDelivery, _ := userChannel2.Consume(
			queue2.Name,
			"user2",
			false,
			true,
			false,
			false,
			nil,
		)
		for delivery := range chanDelivery {
			fmt.Printf("user2 receives message: <%s>\n", string(delivery.Body))
			_ = delivery.Ack(false)
		}
	}()

	select {}
}
