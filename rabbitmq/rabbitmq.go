package rabbitmq

import (
	"context"
	"errors"
	"os"
	"sync"

	"github.com/zikster3262/shared-lib/utils"

	amqp "github.com/rabbitmq/amqp091-go"
)

var (
	ErrNoRabbitMQAddressFound = errors.New("no rabbitMQ address provided")
	rabbitmqLock              sync.Mutex
	rabbitmqLockRW            sync.RWMutex
)

func ConnectToRabbit() (*RabbitMQClient, error) {

	rabbitmqLock.Lock()
	defer rabbitmqLock.Unlock()

	conn, err := amqp.Dial(os.Getenv("RABBITMQ_ADDRESS"))
	utils.FailOnCmpError("rabbitmq", "connection", err)

	utils.LogWithInfo("rabbitmq", "connected to rabbitMQ")

	conn.Close()

	return &RabbitMQClient{
		connection: conn,
	}, err
}

type RabbitMQClient struct {
	connection *amqp.Connection
	channels   []*amqp.Channel
}

func (rmq *RabbitMQClient) CreateChannel() (*amqp.Channel, error) {

	rabbitmqLockRW.Lock()

	chann, err := rmq.connection.Channel()
	if err != nil {
		utils.FailOnCmpError("rabbitmq", "channel", err)
	}
	rmq.channels = append(rmq.channels, chann)

	rabbitmqLockRW.Unlock()

	return chann, err
}

func (rmq *RabbitMQClient) PublishMessage(ctx context.Context, name string, body []byte, ch *amqp.Channel) error {

	rabbitmqLockRW.Lock()

	err := ch.PublishWithContext(ctx,
		"",    // exchange
		name,  // routing key
		false, // mandatory
		false, // immediate
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         body,
			DeliveryMode: amqp.Persistent,
		})
	utils.FailOnCmpError("rabbitmq", "publish", err)

	rabbitmqLockRW.Unlock()
	return err
}

func (rmq *RabbitMQClient) Consume(name string, ch *amqp.Channel) (msgs <-chan amqp.Delivery, err error) {

	rabbitmqLockRW.Lock()

	msgs, err = ch.Consume(
		name,  // queue
		"",    // consumer
		true,  // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // args
	)

	utils.FailOnCmpError("rabbitmq", "consume", err)

	rabbitmqLockRW.Unlock()

	return msgs, err
}
