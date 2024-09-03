package testUtils

import (
	"context"
	"fmt"

	amqpConfig "github.com/nihal-ramaswamy/GoChat/internal/amqp"
	"github.com/nihal-ramaswamy/GoChat/internal/dto"
	"github.com/testcontainers/testcontainers-go/modules/rabbitmq"
)

func GetRabbitMqConfig() *dto.TestConfigDto {
	return &dto.TestConfigDto{
		Username:     "guest",
		Password:     "guest",
		DatabaseName: "go_chat",
	}
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

func SetUpRabbitMqForTesting(ctx context.Context) (*rabbitmq.RabbitMQContainer, *amqpConfig.AmqpConfig, error) {
	testConfig := GetRabbitMqConfig()
	container, err := GetRabbitMqContainer(testConfig, ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get rabbitmq container: %s", err)
	}

	amqpURL, err := container.AmqpURL(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get amqp url: %s", err)
	}

	amqpConfigDto, err := amqpConfig.NewAmqpConfig(amqpURL)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get amqp config: %s", err)
	}

	return container, amqpConfigDto, nil
}
