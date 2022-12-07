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

const (
	timeout = 5
)

type Client struct {
	conn *amqp.Connection
}

func CreateClient() *Client {
	conn, err := amqp.Dial(os.Getenv("RABBITMQ_ADDRESS"))
	utils.FailOnError("rabbitmq", err)

	utils.LogWithInfo("rabbitmq", "connected to rabbitMQ")

	return &Client{
		conn: conn,
	}
}

func (rmq *Client) CreateChannel() *amqp.Channel {
	ch, err := rmq.conn.Channel()
	utils.FailOnError("rabbitmq", err)

	return ch
}

func PublishMessage(channel *amqp.Channel, name string, body []byte) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout*time.Second)
	defer cancel()

	err := channel.PublishWithContext(ctx,
		"",    // exchange
		name,  // routing key
		false, // mandatory
		false, // immediate
		amqp.Publishing{
			Headers:         map[string]interface{}{},
			ContentType:     "application/json",
			ContentEncoding: "",
			DeliveryMode:    amqp.Persistent,
			Priority:        0,
			CorrelationId:   "",
			ReplyTo:         "",
			Expiration:      "",
			MessageId:       "",
			Timestamp:       time.Time{},
			Type:            "",
			UserId:          "",
			AppId:           "",
			Body:            body,
		})
	utils.FailOnError("rabbitmq", err)

	return errors.Unwrap(err)
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
			return nil, errors.Unwrap(err)
		} else {
			return msgs, errors.Unwrap(err)
		}
	}
}

func (rmq *Client) Close() {
	if rmq != nil {
		rmq.Close()
	}
}
