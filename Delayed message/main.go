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

// If you want to use delayed messages, you needs to run the commands below.
//
// wget https://github.com/rabbitmq/rabbitmq-delayed-message-exchange/releases/download/v3.12.0/rabbitmq_delayed_message_exchange-3.12.0.ez -O `rabbitmq-plugins --node rabbit@f0afb4fb72c9 directories -s | head -n 1 | awk '{print $4}'`/rabbitmq_delayed_message_exchange-3.12.0.ez
// rabbitmq-plugins enable rabbitmq_delayed_message_exchange
func main() {
	conn, err := amqp.Dial("amqp://goadmin:123456@localhost:5672/goapp-vhost")
	FailedOnError(err, "failed to create connection to RabbitMQ server")
	defer conn.Close()

	ch, err := conn.Channel()
	FailedOnError(err, "Failed to open a channel")
	defer ch.Close()

	err = ch.ExchangeDeclare(
		"goapp.exchange.delay", // Exchange name
		"x-delayed-message",    // Exchange type
		true,                   // Durable
		false,                  // Auto-deleted
		false,                  // Internal
		false,                  // No-wait
		amqp.Table{
			"x-delayed-type": "direct", // Set delayed exchange type as direct
		},
	)
	FailedOnError(err, "Failed to declare the exchange")

	queue, err := ch.QueueDeclare("goapp.queue.delayed", true, false, false, false, nil)
	FailedOnError(err, "Failed to declare queue")

	err = ch.QueueBind(queue.Name, "", "goapp.exchange.delay", false, nil)
	FailedOnError(err, "Failed to bind queue")

	err = ch.PublishWithContext(context.Background(), "goapp.exchange.delay", "", false, false, amqp.Publishing{
		Headers: map[string]interface{}{
			"x-delay": int64(150000),
		},
		ContentType:  "text/plain",
		Body:         []byte("Delayed message"),
		DeliveryMode: amqp.Persistent,
	})
	FailedOnError(err, "Failed to publish a delayed message")

	deliveryChan, err := ch.Consume(queue.Name, "", true, false, false, false, nil)
	go func() {
		fmt.Println(time.Now())
		for d := range deliveryChan {
			fmt.Printf("%s %v", string(d.Body), time.Now())
		}
	}()
	select {}
}
