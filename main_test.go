package main

import (
	"context"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/nihal-ramaswamy/GoChat/internal/constants"
	"github.com/nihal-ramaswamy/GoChat/internal/dto"
	"github.com/nihal-ramaswamy/GoChat/internal/fx_utils"
	"github.com/nihal-ramaswamy/GoChat/internal/routes"
	"github.com/nihal-ramaswamy/GoChat/internal/testUtils"
	"github.com/nihal-ramaswamy/GoChat/internal/utils"
)

// Test /healthcheck/healtcheck
func TestHealthcheck(t *testing.T) {
	ctx := context.Background()

	rootDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Error getting working directory: %s", err)
	}

	postgresContainer, db, err := testUtils.SetUpPostgresForTesting(ctx, rootDir)
	if err != nil {
		t.Fatalf("Error setting up postgres for testing: %s", err)
	}

	t.Cleanup(func() {
		if err := postgresContainer.Terminate(ctx); err != nil {
			t.Fatalf("failed to terminate container: %s", err)
		}
		db.Close()
	})

	rabbitmqContainer, amqpConfig, err := testUtils.SetUpRabbitMqForTesting(ctx)
	if err != nil {
		t.Fatalf("Error setting up rabbitmq for testing: %s", err)
	}

	t.Cleanup(func() {
		if err := rabbitmqContainer.Terminate(ctx); err != nil {
			t.Fatalf("failed to terminate container: %s", err)
		}
	})

	redisContainer, rdb, err := testUtils.SetUpRedisForTesting(ctx)
	if err != nil {
		t.Fatalf("Error setting up redis for testing: %s", err)
	}
	t.Cleanup(func() {
		if err := redisContainer.Terminate(ctx); err != nil {
			t.Fatalf("failed to terminate container: %s", err)
		}
	})

	upgrader := fx_utils.NewWebsocketUpgrader()
	webscoketMap := dto.NewWebsocketConnectionMap()

	os.Setenv(constants.ENV, "test")
	log := utils.NewZapLogger()

	gin.SetMode(gin.TestMode)
	server := gin.Default()

	routes.NewRoutes(server, db, rdb, ctx, log, amqpConfig, upgrader, webscoketMap)
}
