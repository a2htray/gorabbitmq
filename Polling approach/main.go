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

	producerCh, err := conn.Channel()
	FailedOnError(err, "create producer channel failed")
	defer producerCh.Close()

	queue, err := producerCh.QueueDeclare(
		"goapp.queue.test",
		false,
		false,
		false,
		false,
		nil,
	)
	FailedOnError(err, "declare queue failed")

	err = producerCh.ExchangeDeclare(
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
	err = producerCh.QueueBind(queue.Name, routingKey, "goapp.exchange.direct", false, nil)

	go func() {
		cnt := 0
		interval := time.Second * 5
		ticker := time.NewTicker(interval)
		for _ = range ticker.C {
			now := time.Now()
			fmt.Println(now, "<producer> send message to exchange", cnt)
			ctx := context.Background()
			producerCh.PublishWithContext(ctx, "goapp.exchange.direct", routingKey, false, false, amqp.Publishing{
				ContentType:  "plain/text",
				DeliveryMode: 0,
				Body:         []byte(fmt.Sprintf("message at %v", now)),
			})
			cnt++
		}
	}()

	consumerCh, err := conn.Channel()
	FailedOnError(err, "create consumer channel failed")

	go func() {
		interval := time.Second * 2
		ticker := time.NewTicker(interval)
		for _ = range ticker.C {
			now := time.Now()
			fmt.Println(now, "<consumer> try to receive message from queue")
			delivery, ok, err := consumerCh.Get(queue.Name, false)
			if ok {
				fmt.Println(string(delivery.Body))
			} else {
				if err != nil {
					fmt.Println("error: ", err.Error())
				}
				fmt.Println("the queue is empty now")
			}
		}
	}()

	select {}
}
