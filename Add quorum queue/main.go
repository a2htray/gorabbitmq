package main

import (
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	conn, _ := amqp.Dial("amqp://goadmin:123456@localhost:5672/goapp-vhost")
	defer conn.Close()

	channel, _ := conn.Channel()

	// allow to the capacity of coping queues
	queue, _ := channel.QueueDeclare(
		"goapp.queue.quorum",
		true,
		false,
		false,
		false,
		map[string]interface{}{
			"x-queue-type": "quorum",
		},
	)

	fmt.Println(queue)
}
