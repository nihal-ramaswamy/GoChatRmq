package amqpConfig

import (
	"fmt"

	"github.com/nihal-ramaswamy/GoChat/internal/constants"
	"github.com/nihal-ramaswamy/GoChat/internal/utils"
	amqp "github.com/rabbitmq/amqp091-go"
)

type AmqpConfig struct {
	Host    string
	Conn    *amqp.Connection
	Channel *amqp.Channel
	Queue   amqp.Queue
}

func NewAmqpConfig(host string) (*AmqpConfig, error) {
	conn, err := amqp.Dial(host)
	if err != nil {
		panic(fmt.Errorf("Failed to connect to RabbitMQ: %s", err))
	}
	ch, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("Failed to open channel: %s", err)
	}

	err = ch.ExchangeDeclare(
		constants.EXCHANGE_NAME,
		"topic",
		true,
		false,
		false,
		false,
		nil)
	if err != nil {
		return nil, fmt.Errorf("Failed to declare exchange: %s", err)
	}

	queue, err := ch.QueueDeclare(
		"",    // name
		false, // durable
		false, // delete when unused
		true,  // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		panic(fmt.Errorf("Failed to declare queue: %s", err))
	}

	return &AmqpConfig{
		Host:    host,
		Conn:    conn,
		Channel: ch,
		Queue:   queue,
	}, nil
}

func DefaultAmqpConfig() *AmqpConfig {
	// return NewAmqpConfig(utils.GetDotEnvVariable(constants.RABBITMQ_HOST))
	config, err := NewAmqpConfig(utils.GetDotEnvVariable(constants.RABBITMQ_HOST))
	if err != nil {
		panic(fmt.Errorf("Failed to get amqp config: %s", err))
	}
	return config
}
