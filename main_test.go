package main

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
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

func readFromUserDb(db *sql.DB, email string) (int, error) {
	cnt := 0
	query := `SELECT COUNT(*) FROM "USER" WHERE EMAIL = $1`
	err := db.QueryRow(query, email).Scan(&cnt)
	if err != nil {
		return -1, err
	}
	return cnt, err
}

func readFromRedis(rdb *rdb.Client, key string) (bool, string, error) {
	exists := rdb.Exists(context.Background(), key).Val()
	if exists == 0 {
		return false, "", nil
	}
	val, err := rdb.Get(context.Background(), key).Result()
	if err != nil {
		return true, "", err
	}
	return true, val, nil
}

func setUpRouter(rootDir string, ctx context.Context) (*TestConfig, error) {
	postgresContainer, db, err := testUtils.SetUpPostgresForTesting(ctx, rootDir)
	if err != nil {
		return nil, fmt.Errorf("PostgresContainer error: %s", err)
	}

	rabbitmqContainer, amqpConfig, err := testUtils.SetUpRabbitMqForTesting(ctx)
	if err != nil {
		return nil, fmt.Errorf("RabbitmqContainer Error: %s", err)
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

// Test /auth/register
// Tests creating a new user. Expects 202 status
// Tests creating a new user with the same email. Expects 400 status
// Tests creating a new user with same name, different email. Expects 202 status
// Tests if users are being inserted into the database
func TestAuthRegister(t *testing.T) {
	ctx := context.Background()

	rootDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Error getting working directory: %s", err)
	}

	testConfig, err := setUpRouter(rootDir, ctx)
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

	cnt, err := readFromUserDb(testConfig.Db, userDto.Email)
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

	cnt, err = readFromUserDb(testConfig.Db, userDto.Email)
	if err != nil {
		t.Fatalf("Error reading from db: %s", err)
	}
	if cnt != 1 {
		t.Errorf("Expected count 1, got %d", cnt)
	}
}

// Test /auth/signin and /auth/signout
// Tests if key is inserted and deleted appropriately from redis
func TestLoginLogout(t *testing.T) {
	ctx := context.Background()

	rootDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Error getting working directory: %s", err)
	}

	testConfig, err := setUpRouter(rootDir, ctx)
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

	exists, rdb_val, err := readFromRedis(testConfig.Rdb, userDto.Email)
	if err != nil {
		t.Fatalf("Error reading from redis: %s", err)
	}
	if !exists || rdb_val == "" {
		t.Errorf("Expected token, got empty string")
	}

	w = httptest.NewRecorder()
	req, err = http.NewRequest("POST", "/auth/signout", nil)
	if err != nil {
		t.Fatalf("Error creating request: %s", err)
	}
	req.Header.Set("Token", fmt.Sprintf("Bearer %s", token.Token))
	testConfig.Server.ServeHTTP(w, req)
	if w.Code != 202 {
		t.Errorf("Expected status code: 202, got %d", w.Code)
	}

	exists, rdb_val, err = readFromRedis(testConfig.Rdb, userDto.Email)
	if exists {
		t.Errorf("Expected token to be deleted, got %s", rdb_val)
	}
}
