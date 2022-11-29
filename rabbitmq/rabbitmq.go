package rabbitmq

import (
	"context"
	"errors"
	"os"
	"time"

	"github.com/zikster3262/shared-lib/utils"

	"github.com/rabbitmq/amqp091-go"
	amqp "github.com/rabbitmq/amqp091-go"
)

var (
	ErrNoRabbitMQAddressFound = errors.New("no rabbitMQ address provided")
)

func ConnectToRabbit() (*amqp.Connection, error) {

	conn, err := amqp.Dial(os.Getenv("RABBITMQ_ADDRESS"))
	utils.FailOnError("rabbitmq", err)

	// confirms := make(chan amqp.Confirmation)
	// ch.NotifyPublish(confirms)
	// go func() {
	// 	for confirm := range confirms {
	// 		if confirm.Ack {
	// 			utils.LogWithInfo("rabbitmq", "Confirmed")
	// 		} else {
	// 			utils.LogWithInfo("rabbitmq", "Failed")
	// 		}
	// 	}
	// }()

	// err = ch.Confirm(false)
	// utils.FailOnError("rabbitmq", err)

	utils.LogWithInfo("rabbitmq", "connected to rabbitMQ")
	return conn, err
}

type RabbitMQClient struct {
	channel *amqp091.Connection
}

func CreateRabbitMQClient(r *amqp091.Connection) *RabbitMQClient {
	return &RabbitMQClient{
		channel: r,
	}
}

func (rmq *RabbitMQClient) PublishMessage(name string, body []byte) error {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	defer rmq.channel.Close()

	ch, err := rmq.channel.Channel()
	utils.FailOnError("rabbitmq", err)

	defer ch.Close()

	err = ch.PublishWithContext(ctx,
		"",    // exchange
		name,  // routing key
		false, // mandatory
		false, // immediate
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         body,
			DeliveryMode: amqp.Persistent,
		})
	utils.FailOnError("rabbitmq", err)
	return err
}

func (rmq *RabbitMQClient) Consume(name string) (<-chan amqp.Delivery, error) {

	defer rmq.channel.Close()

	ch, err := rmq.channel.Channel()
	utils.FailOnError("rabbitmq", err)

	defer ch.Close()

	msgs, err := ch.Consume(
		name,  // queue
		"",    // consumer
		true,  // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // args
	)

	return msgs, err
}
