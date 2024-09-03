package utils

import (
	"context"
	"database/sql"
	"fmt"
	"math/rand"
	"path/filepath"
	"time"

	"github.com/nihal-ramaswamy/GoChat/internal/constants"
	"github.com/nihal-ramaswamy/GoChat/internal/dto"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/modules/rabbitmq"
	"github.com/testcontainers/testcontainers-go/wait"
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
		return nil, fmt.Errorf("Failed to declare queue: %s", err)
	}

	return &AmqpConfig{
		Host:    host,
		Conn:    conn,
		Channel: ch,
		Queue:   queue,
	}, nil
}

func GetDbConfig() *dto.TestConfigDto {
	return &dto.TestConfigDto{
		Username:     "postgresTest",
		Password:     "postgresTest",
		DatabaseName: "go_chat",
	}
}

func GetRabbitMqConfig() *dto.TestConfigDto {
	return &dto.TestConfigDto{
		Username:     "guest",
		Password:     "guest",
		DatabaseName: "go_chat",
	}
}

func GetPostgresContainer(
	testConfig *dto.TestConfigDto,
	rootDir string,
	ctx context.Context,
) (*postgres.PostgresContainer, error) {
	container, err := postgres.Run(
		ctx,
		"docker.io/postgres:16-alpine",
		postgres.WithInitScripts(filepath.Join(rootDir, "db", "init.sql")),
		postgres.WithUsername(testConfig.Username),
		postgres.WithPassword(testConfig.Password),
		postgres.WithDatabase(testConfig.DatabaseName),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).WithStartupTimeout(5*time.Second)),
		// postgres.WithSQLDriver("pq"),
	)

	return container, err
}

func GetRabbitMqContainer(
	testConfig *dto.TestConfigDto,
	ctx context.Context,
) (*rabbitmq.RabbitMQContainer, error) {
	rabbitmqContainer, err := rabbitmq.Run(ctx,
		"rabbitmq:3.12.11-management-alpine",
		rabbitmq.WithAdminUsername(testConfig.Username),
		rabbitmq.WithAdminPassword(testConfig.Password),
	)

	return rabbitmqContainer, err
}

func RandStringRunes(n int) string {
	letterRunes := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func FilterChatsByUserId(chats []*dto.Chat, userId string) int {
	count := 0
	for _, chat := range chats {
		if chat.SenderId == userId || chat.ReceiverId == userId {
			count++
		}
	}
	return count
}

func SetUpPostgresForTesting(ctx context.Context, rootDir string) (*postgres.PostgresContainer, *sql.DB, error) {
	testConfig := GetDbConfig()

	container, err := GetPostgresContainer(testConfig, rootDir, ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get postgres container: %s", err)
	}

	dbURL, err := container.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get connection string: %s", err)
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to open connection: %s", err)
	}

	err = db.Ping()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to ping db: %s", err)
	}

	return container, db, nil
}

func SetUpRabbitMqForTesting(ctx context.Context) (*rabbitmq.RabbitMQContainer, *AmqpConfig, error) {
	testConfig := GetRabbitMqConfig()
	container, err := GetRabbitMqContainer(testConfig, ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get rabbitmq container: %s", err)
	}

	amqpURL, err := container.AmqpURL(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get amqp url: %s", err)
	}

	amqpConfigDto, err := NewAmqpConfig(amqpURL)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get amqp config: %s", err)
	}

	return container, amqpConfigDto, nil
}
