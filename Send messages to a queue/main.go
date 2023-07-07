package main

import (
	"context"
	"encoding/json"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"time"
)

func main() {
	// 1. create connection
	conn, err := amqp.Dial("amqp://goadmin:123456@localhost:5672/goapp-vhost")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	// 2. create channel
	channel, err := conn.Channel()
	if err != nil {
		log.Fatal(err)
	}
	defer channel.Close()

	// 3. declare queue
	queue, err := channel.QueueDeclare(
		"goapp.queue.test",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatal(err)
	}

	// 4. declare exchange
	err = channel.ExchangeDeclare(
		"goapp.exchange.direct",
		amqp.ExchangeDirect,
		true,
		false,
		false,
		true,
		nil,
	)
	if err != nil {
		log.Fatal(err)
	}

	// 5. bind queue to exchange, and declare routing key
	routingKey := "routing.key"
	err = channel.QueueBind(queue.Name, routingKey, "goapp.exchange.direct", false, nil)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	for i := 0; i < 10; i++ {
		data := map[string]any{
			"id":       i,
			"username": fmt.Sprintf("user%d", i),
		}

		b, err := json.Marshal(data)
		if err != nil {
			log.Fatal(err)
		}

		// 6. specify the exchange and the routing key and send messages to the exchange
		err = channel.PublishWithContext(
			ctx,
			"goapp.exchange.direct",
			routingKey,
			false,
			false,
			amqp.Publishing{
				DeliveryMode: amqp.Persistent,
				ContentType:  "application/json",
				Body:         b,
			})

		if err != nil {
			log.Fatal(err)
		}
	}

	//select {}
}
