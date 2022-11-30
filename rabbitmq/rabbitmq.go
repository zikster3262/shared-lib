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

	conn, err := amqp.Dial(os.Getenv("RABBITMQ_ADDRESS"))
	utils.FailOnError("rabbitmq", err)

	ch, err := conn.Channel()
	utils.FailOnError("rabbitmq", err)

	utils.LogWithInfo("rabbitmq", "connected to rabbitMQ")
	return ch, err
}

type RabbitMQClient struct {
	channel *amqp.Channel
}

func CreateRabbitMQClient(r *amqp.Channel) *RabbitMQClient {
	return &RabbitMQClient{
		channel: r,
	}
}

func (rmq *RabbitMQClient) PublishMessage(ctx context.Context, name string, body []byte) error {

	err := rmq.channel.PublishWithContext(ctx,
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

	msgs, err := rmq.channel.Consume(
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

func (rmq *RabbitMQClient) Close() {

	if rmq != nil {
		rmq.Close()
	}
}
