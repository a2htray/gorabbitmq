package main

import (
	"context"
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
	_ = channel.ExchangeDeclare(
		"goapp.exchange.headers.01",
		amqp.ExchangeHeaders,
		true,
		false,
		false,
		false,
		nil,
	)

	returnChan := make(chan amqp.Return)
	channel.NotifyReturn(returnChan)

	go func() {
		count := 0
		for v := range returnChan {
			fmt.Println(count, string(v.Body))
			count++
		}
	}()

	queue, _ := channel.QueueDeclare(
		"goapp.queue.headers.01",
		true,
		false,
		false,
		false,
		nil,
	)

	_ = channel.QueueBind(
		queue.Name,
		"",
		"goapp.exchange.headers.01",
		false,
		amqp.Table(map[string]interface{}{
			"x-match":       "all", // all or any
			"apiVersion":    "0.0.1",
			"clientVersion": "0.0.1a",
		}),
	)

	go func() {
		for i := 0; i < 5; i++ {
			channel.PublishWithContext(
				context.Background(),
				"goapp.exchange.headers.01",
				"",
				true,
				false,
				amqp.Publishing{
					ContentType: "plain/text",
					Body:        []byte(fmt.Sprintf("without header properties %d\n", i)),
				},
			)
		}
	}()

	go func() {
		for i := 0; i < 5; i++ {
			channel.PublishWithContext(
				context.Background(),
				"goapp.exchange.headers.01",
				"",
				false,
				false,
				amqp.Publishing{
					ContentType: "plain/text",
					Headers: map[string]interface{}{
						"apiVersion":    "0.0.1",
						"clientVersion": "0.0.1a",
					},
					Body: []byte(fmt.Sprintf("with header properties %d\n", i)),
				},
			)
		}
	}()

	select {}

}
