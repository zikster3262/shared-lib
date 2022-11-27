package rabbitmq

import (
	"context"
	"errors"
	"os"

	"github.com/zikster3262/shared-lib/utils"

	amqp "github.com/rabbitmq/amqp091-go"
)

var (
	ErrNoRabbitMQAddressFound = errors.New("no rabbitMQ address provided")
)

func ConnectToRabbit() (*amqp.Channel, error) {

	addr := os.Getenv("RABBITMQ_ADDRESS")

	if addr == "" {
		utils.FailOnError("rabbitmq", ErrNoRabbitMQAddressFound)
	}

	conn, err := amqp.Dial(addr)
	utils.FailOnError("rabbitmq", err)

	ch, err := conn.Channel()
	utils.FailOnError("rabbitmq", err)

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

func (rmq *RabbitMQClient) PublishMessage(name string, ctx context.Context, body []byte) error {

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
