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

type ExchangeConfig struct {
	Name       string
	Kind       string
	Durable    bool
	AutoDelete bool
	Internal   bool
	NoWait     bool
	Args       amqp.Table
}

type QueueConfig struct {
	Name       string
	Durable    bool
	AutoDelete bool
	Exclusive  bool
	NoWait     bool
	Args       amqp.Table
	RoutingKey string
}

type Producer struct {
	conn           *amqp.Connection
	ch             *amqp.Channel
	queue          *amqp.Queue
	exchangeConfig ExchangeConfig
	queueConfig    QueueConfig
}

func (p *Producer) BindExAndQ(exchangeConfig ExchangeConfig, queueConfig QueueConfig) error {
	if err := p.ch.ExchangeDeclare(
		exchangeConfig.Name,
		exchangeConfig.Kind,
		exchangeConfig.Durable,
		exchangeConfig.AutoDelete,
		exchangeConfig.Internal,
		exchangeConfig.NoWait,
		exchangeConfig.Args,
	); err != nil {
		return err
	}

	if queue, err := p.ch.QueueDeclare(
		queueConfig.Name,
		queueConfig.Durable,
		queueConfig.AutoDelete,
		queueConfig.Exclusive,
		queueConfig.NoWait,
		queueConfig.Args,
	); err != nil {
		return err
	} else {
		if err := p.ch.QueueBind(
			queue.Name,
			queueConfig.RoutingKey,
			exchangeConfig.Name,
			false,
			nil,
		); err != nil {
			return err
		}
	}
	p.exchangeConfig = exchangeConfig
	p.queueConfig = queueConfig
	return nil
}

func (p *Producer) Publish(b []byte) error {
	return p.ch.PublishWithContext(
		context.TODO(),
		p.exchangeConfig.Name,
		p.queueConfig.RoutingKey,
		false,
		false,
		amqp.Publishing{
			DeliveryMode: 2,
			ContentType:  "plain/text",
			Body:         b,
		},
	)
}

func NewProducer(conn *amqp.Connection) *Producer {
	producer := &Producer{
		conn: conn,
	}
	ch, err := producer.conn.Channel()
	FailedOnError(err, "create producer channel failed")
	producer.ch = ch
	return producer
}

type Consumer struct {
	conn *amqp.Connection
	ch   *amqp.Channel
}

func (c *Consumer) Subscribe(queueName string, callback func([]byte)) error {
	ch, err := c.ch.Consume(queueName, "", true, false, false, false, nil)
	if err != nil {
		return err
	}
	for delivery := range ch {
		callback(delivery.Body)
	}
	return nil
}

func NewConsumer(conn *amqp.Connection) *Consumer {
	var err error
	consumer := &Consumer{conn: conn}
	consumer.ch, err = consumer.conn.Channel()
	FailedOnError(err, "create consumer channel failed")
	return consumer
}

func ConsumeMessage(b []byte) {
	fmt.Println(string(b))
}

func main() {
	conn, err := amqp.Dial("amqp://goadmin:123456@localhost:5672/goapp-vhost")
	FailedOnError(err, "could not connect to RabbitMQ server")

	exchangeConfig := ExchangeConfig{
		Name:       "goapp.exchange.direct",
		Kind:       amqp.ExchangeDirect,
		Durable:    true,
		AutoDelete: false,
		Internal:   false,
		NoWait:     false,
		Args:       nil,
	}
	queueConfig := QueueConfig{
		Name:       "goapp.queue.test",
		Durable:    false,
		AutoDelete: false,
		Exclusive:  false,
		NoWait:     false,
		Args:       nil,
		RoutingKey: "routing.key",
	}

	producer := NewProducer(conn)
	err = producer.BindExAndQ(exchangeConfig, queueConfig)
	FailedOnError(err, "bind queue to exchange failed")

	consumer := NewConsumer(conn)

	go func() {
		ticker := time.NewTicker(time.Second * 3)
		for _ = range ticker.C {
			fmt.Println("try to send message to RabbitMQ server")
			if err = producer.Publish([]byte(fmt.Sprintf("%v", time.Now()))); err != nil {
				FailedOnError(err, "publish message failed")
			}

		}
	}()

	go func() {
		if err = consumer.Subscribe(queueConfig.Name, ConsumeMessage); err != nil {
			FailedOnError(err, "consume message failed")
		}
	}()

	select {}
}
