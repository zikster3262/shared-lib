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

type RabbitMQClient struct {
	conn *amqp.Connection
}

func CreateRabbitMQClient() *RabbitMQClient {

	conn, err := amqp.Dial(os.Getenv("RABBITMQ_ADDRESS"))
	utils.FailOnError("rabbitmq", err)

	utils.LogWithInfo("rabbitmq", "connected to rabbitMQ")
	return &RabbitMQClient{
		conn: conn,
	}
}

func (r *RabbitMQClient) CreateChannel() *amqp.Channel {
	ch, err := r.conn.Channel()
	utils.FailOnError("rabbitmq", err)
	return ch
}

func PublishMessage(channel *amqp.Channel, name string, body []byte) error {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := channel.PublishWithContext(ctx,
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

func Consume(channel *amqp.Channel, name string) (<-chan amqp.Delivery, error) {

	msgs, err := channel.Consume(
		name,  // queue
		"",    // consumer
		true,  // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // args
	)

	for {
		_, ok := <-msgs
		if !ok {
			return nil, nil
		} else {
			return msgs, err
		}
	}

}

func (rmq *RabbitMQClient) Close() {

	if rmq != nil {
		rmq.Close()
	}
}
