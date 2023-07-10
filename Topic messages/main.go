package main

import (
	"context"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"math/rand"
	"time"
)

func FailedOnError(err error, message string) {
	if err != nil {
		log.Fatal(message, "Error: ", err.Error())
	}
}

func CreateTopicExchange(conn *amqp.Connection, name string) {
	channel, err := conn.Channel()
	FailedOnError(err, "<CreateTopicExchange> create channel failed")
	defer channel.Close()

	err = channel.ExchangeDeclare(
		name,
		amqp.ExchangeTopic,
		true,
		true,
		false,
		false,
		nil,
	)
	FailedOnError(err, "<CreateTopicExchange> create exchange failed")
}

func BindExchangeAndQueue(conn *amqp.Connection, queueName, exchangeName, routingKey string) {
	channel, err := conn.Channel()
	FailedOnError(err, "<BindExchangeAndQueue> create channel failed")
	defer channel.Close()

	err = channel.QueueBind(
		queueName,
		routingKey,
		exchangeName,
		false,
		nil,
	)
	FailedOnError(err, "<BindExchangeAndQueue> bind queue failed")
}

func CreateQueue(conn *amqp.Connection, name string) amqp.Queue {
	channel, err := conn.Channel()
	FailedOnError(err, "<CreateQueue> create channel failed")
	defer channel.Close()

	queue, err := channel.QueueDeclare(
		name,
		true,
		false,
		false,
		false,
		nil,
	)
	FailedOnError(err, "<CreateQueue> create queue failed")

	return queue
}

func Send(conn *amqp.Connection, exchangeName string) {
	channel, err := conn.Channel()
	FailedOnError(err, "<Send> create channel failed")
	defer channel.Close()

	routingKeys := []string{"apple.123", "banana.123"}
	ctx := context.Background()
	for i := 0; i < 10; i++ {
		idx := rand.Intn(2)
		routingKey := routingKeys[idx]
		err = channel.PublishWithContext(
			ctx,
			exchangeName,
			routingKey,
			false,
			false,
			amqp.Publishing{
				ContentType:  "plain/text",
				DeliveryMode: 2,
				Body:         []byte(fmt.Sprintf("%v", time.Now())),
			},
		)
		fmt.Println(fmt.Sprintf("<Send> send message with routing key %s to exchange %s", routingKey, exchangeName))
	}
	FailedOnError(err, "<Send> send message failed")
}

func Receive(conn *amqp.Connection, consumerName, queueName string, callback func(string, amqp.Delivery)) {
	fmt.Printf("<Receive> consumer %s starts to receive messages\n", consumerName)
	channel, err := conn.Channel()
	FailedOnError(err, "<Receive> create channel failed")
	defer channel.Close()

	chDelivery, err := channel.Consume(
		queueName,
		consumerName,
		false,
		false,
		false,
		false,
		nil,
	)
	FailedOnError(err, fmt.Sprintf("<Receive> consumer %s fails to receive messages", consumerName))
	for delivery := range chDelivery {
		callback(consumerName, delivery)
	}
}

func main() {
	conn, err := amqp.Dial("amqp://goadmin:123456@localhost:5672/goapp-vhost")
	FailedOnError(err, "connect to RabbitMQ server failed")
	defer conn.Close()

	topicExchangeName := "goapp.exchange.topic.product"
	CreateTopicExchange(conn, topicExchangeName)

	aQueue := CreateQueue(conn, "app.queue.a")
	bQueue := CreateQueue(conn, "app.queue.b")

	BindExchangeAndQueue(conn, aQueue.Name, topicExchangeName, "apple.*")
	BindExchangeAndQueue(conn, bQueue.Name, topicExchangeName, "banana.*")

	go func() {
		Send(conn, topicExchangeName)
	}()

	var callback = func(consumerName string, delivery amqp.Delivery) {
		fmt.Println(consumerName, string(delivery.Body), delivery.RoutingKey)
	}

	// consume messages with topic 'apple.*'
	go func() {
		Receive(conn, "consumer.a", aQueue.Name, callback)
	}()

	// consume messages with topic 'banana.*'
	go func() {
		Receive(conn, "consumer.b", bQueue.Name, callback)
	}()

	select {}
}
