package main

import (
	"context"
	"database/sql"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	_ "github.com/lib/pq"
	amqpConfig "github.com/nihal-ramaswamy/GoChat/internal/amqp"
	"github.com/nihal-ramaswamy/GoChat/internal/constants"
	"github.com/nihal-ramaswamy/GoChat/internal/dto"
	"github.com/nihal-ramaswamy/GoChat/internal/fx_utils"
	"github.com/nihal-ramaswamy/GoChat/internal/routes"
	"github.com/nihal-ramaswamy/GoChat/internal/testUtils"
	"github.com/nihal-ramaswamy/GoChat/internal/utils"
	rdb "github.com/redis/go-redis/v9"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/modules/rabbitmq"
	"github.com/testcontainers/testcontainers-go/modules/redis"
	"go.uber.org/zap"
)

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

func setUpRouter(rootDir string, ctx context.Context) (*TestConfig, error) {
	postgresContainer, db, err := testUtils.SetUpPostgresForTesting(ctx, rootDir)
	if err != nil {
		return nil, err
	}

	rabbitmqContainer, amqpConfig, err := testUtils.SetUpRabbitMqForTesting(ctx)
	if err != nil {
		return nil, err
	}

	redisContainer, rdb, err := testUtils.SetUpRedisForTesting(ctx)
	upgrader := fx_utils.NewWebsocketUpgrader()
	webscoketMap := dto.NewWebsocketConnectionMap()

	os.Setenv(constants.ENV, "test")
	log := utils.NewZapLogger()

	gin.SetMode(gin.TestMode)
	server := gin.Default()

	return &TestConfig{
		PostgresContainer: postgresContainer,
		Db:                db,
		RabbitmqContainer: rabbitmqContainer,
		AmqpConfig:        amqpConfig,
		RedisContainer:    redisContainer,
		Rdb:               rdb,
		Server:            server,
		Log:               log,
		Upgrader:          upgrader,
		WebsocketMap:      webscoketMap,
	}, nil
}

// Test /healthcheck/healtcheck
func TestHealthcheck(t *testing.T) {
	ctx := context.Background()

	rootDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Error getting working directory: %s", err)
	}

	testConfig, err := setUpRouter(rootDir, ctx)
	if err != nil {
		t.Fatalf("Error setting up postgres for testing: %s", err)
	}

	t.Cleanup(func() {
		if err := testConfig.PostgresContainer.Terminate(ctx); err != nil {
			t.Fatalf("failed to terminate container: %s", err)
		}
		testConfig.Db.Close()
	})

	if err != nil {
		t.Fatalf("Error setting up rabbitmq for testing: %s", err)
	}

	t.Cleanup(func() {
		if err := testConfig.RabbitmqContainer.Terminate(ctx); err != nil {
			t.Fatalf("failed to terminate container: %s", err)
		}
	})

	if err != nil {
		t.Fatalf("Error setting up redis for testing: %s", err)
	}
	t.Cleanup(func() {
		if err := testConfig.RedisContainer.Terminate(ctx); err != nil {
			t.Fatalf("failed to terminate container: %s", err)
		}
	})

	routes.NewRoutes(
		testConfig.Server,
		testConfig.Db,
		testConfig.Rdb,
		ctx,
		testConfig.Log,
		testConfig.AmqpConfig,
		testConfig.Upgrader,
		testConfig.WebsocketMap,
	)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/healthcheck/healthcheck", nil)
	testConfig.Server.ServeHTTP(w, req)
	if w.Code != 200 {
		t.Errorf("Expected status code 200, got %d", w.Code)
	}
}
