package middlewares_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	_ "github.com/lib/pq"
	"github.com/nihal-ramaswamy/GoChat/internal/dto"
	"github.com/nihal-ramaswamy/GoChat/internal/routes"
	"github.com/nihal-ramaswamy/GoChat/internal/testUtils"
)

// Test auth middleware
func TestAuthMiddleware(t *testing.T) {
	ctx := context.Background()

	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Error getting working directory: %s", err)
	}
	rootDir := filepath.Join(wd, "..", "..")

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

	// Rgister user
	w := httptest.NewRecorder()
	userDto := dto.User{
		Name:     "test",
		Password: "test",
		Email:    "test",
	}
	userJson, err := json.Marshal(userDto)
	if err != nil {
		t.Fatalf("Error converting userDto to json: %s", err)
	}
	req, err := http.NewRequest("POST", "/auth/register", bytes.NewBuffer(userJson))
	if err != nil {
		t.Fatalf("Error creating request: %s", err)
	}

	testConfig.Server.ServeHTTP(w, req)
	if w.Code != 202 {
		t.Errorf("Expected status code: 202, got %d", w.Code)
	}
	body, err := io.ReadAll(w.Body)
	if err != nil {
		t.Fatalf("Error reading body: %s", err)
	}
	var result testUtils.IdDto
	err = json.Unmarshal(body, &result)
	if err != nil {
		t.Fatalf("Error unmarshalling response: %s", err)
	}
	if result.Id == "" {
		t.Errorf("Expected id, got empty string")
	}

	// Sign in user
	w = httptest.NewRecorder()
	req, err = http.NewRequest("POST", "/auth/signin", bytes.NewBuffer(userJson))
	if err != nil {
		t.Fatalf("Error creating request: %s", err)
	}
	testConfig.Server.ServeHTTP(w, req)
	if w.Code != 202 {
		t.Errorf("Expected status code: 202, got %d", w.Code)
	}
	body, err = io.ReadAll(w.Body)
	if err != nil {
		t.Fatalf("Error reading body: %s", err)
	}

	var token testUtils.TokenDto
	err = json.Unmarshal(body, &token)
	if err != nil {
		t.Fatalf("Error unmarshalling response: %s", err)
	}
	if token.Token == "" {
		t.Errorf("Expected token, got empty string")
	}

	// Test auth middleware
	w = httptest.NewRecorder()
	req, err = http.NewRequest("GET", "/healthcheck/healthcheckAuth", nil)
	if err != nil {
		t.Fatalf("Error sending request to healthcheckAuth: %s", err)
	}
	req.Header.Set("Token", fmt.Sprintf("Bearer %s", token.Token))
	testConfig.Server.ServeHTTP(w, req)
	if w.Code != 200 {
		t.Errorf("Expected status code: 200, got %d", w.Code)
	}

	w = httptest.NewRecorder()
	req, err = http.NewRequest("GET", "/healthcheck/healthcheckAuth", nil)
	if err != nil {
		t.Fatalf("Error sending request to healthcheckAuth: %s", err)
	}
	testConfig.Server.ServeHTTP(w, req)
	if w.Code != 401 {
		t.Errorf("Expected status code: 200, got %d", w.Code)
	}
}
