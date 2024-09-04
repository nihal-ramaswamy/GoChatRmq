package testUtils

import (
	"database/sql"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	amqpConfig "github.com/nihal-ramaswamy/GoChat/internal/amqp"
	"github.com/nihal-ramaswamy/GoChat/internal/dto"
	rdb "github.com/redis/go-redis/v9"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/modules/rabbitmq"
	"github.com/testcontainers/testcontainers-go/modules/redis"
	"go.uber.org/zap"
)

type IdDto struct {
	Id string `json:"id"`
}

type ErrorDto struct {
	Error string `json:"error"`
}

type TokenDto struct {
	Token string `json:"token"`
}

type TestConfig struct {
	PostgresContainer *postgres.PostgresContainer
	Db                *sql.DB
	RabbitmqContainer *rabbitmq.RabbitMQContainer
	AmqpConfig        *amqpConfig.AmqpConfig
	RedisContainer    *redis.RedisContainer
	Rdb               *rdb.Client
	Server            *gin.Engine
	Log               *zap.Logger
	Upgrader          *websocket.Upgrader
	WebsocketMap      *dto.WebsocketConnectionMap
}
