package rabbitmq

import (
	"context"
	"errors"
	"os"
	"sync"
	"time"

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
	utils.FailOnError("rabbitmq", err)

	utils.LogWithInfo("rabbitmq", "connected to rabbitMQ")

	return &RabbitMQClient{
		connection: conn,
	}, err
}

type RabbitMQClient struct {
	connection *amqp.Connection
}

func (rmq *RabbitMQClient) PublishMessage(name string, body []byte) error {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rabbitmqLockRW.Lock()

	chann, err := rmq.connection.Channel()
	if err != nil {
		return err
	}

	err = chann.PublishWithContext(ctx,
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

	rabbitmqLockRW.Unlock()
	return err
}

func (rmq *RabbitMQClient) Consume(name string) (<-chan amqp.Delivery, error) {

	rabbitmqLockRW.Lock()

	chann, err := rmq.connection.Channel()
	if err != nil {
		return nil, err
	}

	msgs, err := chann.Consume(
		name,  // queue
		"",    // consumer
		true,  // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // args
	)

	rabbitmqLockRW.Unlock()

	return msgs, err
}
