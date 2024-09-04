package auth_api_test

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/nihal-ramaswamy/GoChat/internal/dto"
	"github.com/nihal-ramaswamy/GoChat/internal/routes"
	"github.com/nihal-ramaswamy/GoChat/internal/testUtils"
)

// Test /auth/register
// Tests creating a new user. Expects 202 status
// Tests creating a new user with the same email. Expects 400 status
// Tests creating a new user with same name, different email. Expects 202 status
// Tests if users are being inserted into the database
func TestAuthRegister(t *testing.T) {
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

	cnt, err := testUtils.ReadFromUserDb(testConfig.Db, userDto.Email)
	if err != nil {
		t.Fatalf("Error reading from db: %s", err)
	}
	if cnt != 1 {
		t.Errorf("Expected count 1, got %d", cnt)
	}

	w = httptest.NewRecorder()
	req, err = http.NewRequest("POST", "/auth/register", bytes.NewBuffer(userJson))
	if err != nil {
		t.Fatalf("Error creating request: %s", err)
	}
	testConfig.Server.ServeHTTP(w, req)
	if w.Code != 400 {
		t.Errorf("Expected status code: 400, got %d", w.Code)
	}
	var resultError testUtils.ErrorDto
	body, err = io.ReadAll(w.Body)
	if err != nil {
		t.Fatalf("Error reading body: %s", err)
	}

	err = json.Unmarshal(body, &resultError)
	if err != nil {
		t.Fatalf("Error unmarshalling response: %s", err)
	}
	if resultError.Error != "User with email test already exists" {
		t.Errorf("Unexpected error message. Expected: User with email test already exists, got: %s", resultError.Error)
	}

	userDto.Email = "test1"
	userJson, err = json.Marshal(userDto)
	if err != nil {
		t.Fatalf("Error converting userDto to json: %s", err)
	}
	req, err = http.NewRequest("POST", "/auth/register", bytes.NewBuffer(userJson))
	if err != nil {
		t.Fatalf("Error creating request: %s", err)
	}

	w1 := httptest.NewRecorder()

	testConfig.Server.ServeHTTP(w1, req)
	if w1.Code != 202 {
		t.Errorf("Expected status code: 202, got %d", w.Code)
	}
	body, err = io.ReadAll(w1.Body)
	if err != nil {
		t.Fatalf("Error reading body: %s", err)
	}
	err = json.Unmarshal(body, &result)
	if err != nil {
		t.Fatalf("Error unmarshalling response: %s", err)
	}
	if result.Id == "" {
		t.Errorf("Expected id, got empty string")
	}

	cnt, err = testUtils.ReadFromUserDb(testConfig.Db, userDto.Email)
	if err != nil {
		t.Fatalf("Error reading from db: %s", err)
	}
	if cnt != 1 {
		t.Errorf("Expected count 1, got %d", cnt)
	}
}
