package main

import (
	"context"
	"os"
	"testing"

	_ "github.com/lib/pq"
	"github.com/nihal-ramaswamy/GoChat/internal/utils"
)

// Test /healthcheck/healtcheck
func TestHealthcheck(t *testing.T) {
	ctx := context.Background()

	rootDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Error getting working directory: %s", err)
	}

	postgresContainer, db, err := utils.SetUpPostgresForTesting(ctx, rootDir)
	if err != nil {
		t.Fatalf("Error setting up postgres for testing: %s", err)
	}

	t.Cleanup(func() {
		if err := postgresContainer.Terminate(ctx); err != nil {
			t.Fatalf("failed to terminate container: %s", err)
		}
		db.Close()
	})

	rabbitmqContainer, _, err := utils.SetUpRabbitMqForTesting(ctx)
	if err != nil {
		t.Fatalf("Error setting up rabbitmq for testing: %s", err)
	}

	t.Cleanup(func() {
		if err := rabbitmqContainer.Terminate(ctx); err != nil {
			t.Fatalf("failed to terminate container: %s", err)
		}
	})

	redisContainer, _, err := utils.SetUpRedisForTesting(ctx)
	if err != nil {
		t.Fatalf("Error setting up redis for testing: %s", err)
	}
	t.Cleanup(func() {
		if err := redisContainer.Terminate(ctx); err != nil {
			t.Fatalf("failed to terminate container: %s", err)
		}
	})
}
