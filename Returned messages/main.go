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

// Returned message is the message that can not be routed
func main() {
	conn, err := amqp.Dial("amqp://goadmin:123456@localhost:5672/goapp-vhost")
	FailedOnError(err, "failed to create connection to RabbitMQ server")
	defer conn.Close()

	channel, _ := conn.Channel()
	returnChan := make(chan amqp.Return)
	channel.NotifyReturn(returnChan)

	go func() {
		count := 0
		for v := range returnChan {
			fmt.Println(count, string(v.Body))
			count++
		}
	}()

	channel.PublishWithContext(
		context.Background(),
		"",
		"goapp.queue.not_exist",
		true,
		false,
		amqp.Publishing{
			ContentType: "plain/text",
			Body:        []byte(fmt.Sprintf("%v", time.Now())),
		},
	)

	select {}

}
