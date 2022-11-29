package rabbitmq

import (
	"context"
	"errors"
	"os"
	"time"

	"github.com/zikster3262/shared-lib/utils"

	amqp "github.com/rabbitmq/amqp091-go"
)

var (
	ErrNoRabbitMQAddressFound = errors.New("no rabbitMQ address provided")
)

func ConnectToRabbit() (*amqp.Channel, error) {

	conn, err := amqp.Dial(os.Getenv("RABBITMQ_ADDRESS"))
	utils.FailOnError("rabbitmq", err)

	ch, err := conn.Channel()
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
	return ch, err
}

type RabbitMQClient struct {
	ch *amqp.Channel
}

func CreateRabbitMQClient(r *amqp.Channel) *RabbitMQClient {
	return &RabbitMQClient{
		ch: r,
	}
}

func (rmq *RabbitMQClient) PublishMessage(name string, body []byte) error {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := rmq.ch.PublishWithContext(ctx,
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

	msgs, err := rmq.ch.Consume(
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
