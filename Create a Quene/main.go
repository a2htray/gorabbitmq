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

	channel, err := conn.Channel()
	if err != nil {
		log.Fatal(err)
	}
	defer channel.Close()

	// channel.QueueDeclare creates a queue if it doesn't exists, or ensures
	// that an existing queue matches the same properties.

	// a default binding to empty exchange "" which type is "direct" will be
	// created
	queue, err := channel.QueueDeclare(
		// queue name, if the queue name is empty, the server will generate a unique name
		"goapp.queue.test",
		// the durable and autoDelete parameters are used together, here are 4 situations
		// 1. durable is true and autoDelete is false: the queues will survive on server restart and remain
		// when there are no remaining consumers or bindings

		// 2. durable is false and autoDelete is true: queues will not be redeclared on server restart and will be
		// deleted when the last consumer is canceled or the consumer's channel is closed

		// 3. durable is false and autoDelete is false: queues will remain declared as long as the server is running

		// 4. durable is true and autoDelete is true: the queues will be restored on server restart. And the queues
		// will be removed without active consumers
		false,
		false,
		// whether to be accessible by the other connections and channels
		false,
		// assume the queue is declared on the server when noWait if true
		false,
		nil,
	)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(queue)
}
