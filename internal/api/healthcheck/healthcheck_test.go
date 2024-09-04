package healthcheck_api_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	_ "github.com/lib/pq"
	"github.com/nihal-ramaswamy/GoChat/internal/routes"
	"github.com/nihal-ramaswamy/GoChat/internal/testUtils"
)

func TestHealthcheck(t *testing.T) {
	ctx := context.Background()

	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Error getting working directory: %s", err)
	}
	rootDir := filepath.Join(wd, "..", "..", "..")

	testConfig, err := testUtils.SetUpRouter(rootDir, ctx)
	if err != nil {
		t.Fatalf("Error setting up server for testing: %s", err)
	}

	t.Cleanup(func() {
		if err := testConfig.PostgresContainer.Terminate(ctx); err != nil {
			t.Fatalf("failed to terminate container: %s", err)
		}
		testConfig.Db.Close()
	})

	t.Cleanup(func() {
		if err := testConfig.RabbitmqContainer.Terminate(ctx); err != nil {
			t.Fatalf("failed to terminate container: %s", err)
		}
	})

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
	req, err := http.NewRequest("GET", "/healthcheck/healthcheck", nil)
	if err != nil {
		t.Fatalf("Error creating request: %s", err)
	}
	testConfig.Server.ServeHTTP(w, req)
	if w.Code != 200 {
		t.Errorf("Expected status code 200, got %d", w.Code)
	}
}
