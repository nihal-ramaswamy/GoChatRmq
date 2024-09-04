package testUtils

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"os"
	"regexp"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/nihal-ramaswamy/GoChat/internal/constants"
	"github.com/nihal-ramaswamy/GoChat/internal/dto"
	"github.com/nihal-ramaswamy/GoChat/internal/fx_utils"
	"github.com/nihal-ramaswamy/GoChat/internal/utils"
)

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

func SetUpRouter(rootDir string, ctx context.Context) (*TestConfig, error) {
	postgresContainer, db, err := SetUpPostgresForTesting(ctx, rootDir)
	if err != nil {
		return nil, fmt.Errorf("PostgresContainer error: %s", err)
	}

	rabbitmqContainer, amqpConfig, err := SetUpRabbitMqForTesting(ctx)
	if err != nil {
		return nil, fmt.Errorf("RabbitmqContainer Error: %s", err)
	}

	redisContainer, rdb, err := SetUpRedisForTesting(ctx)
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

// LoadEnv loads env vars from .env
func LoadEnv() {
	projectDirName := "go_chat"
	re := regexp.MustCompile(`^(.*` + projectDirName + `)`)
	cwd, _ := os.Getwd()
	rootPath := re.Find([]byte(cwd))

	err := godotenv.Load(string(rootPath) + `/.env`)
	if err != nil {
		log.Fatalln("Problem loading .env file")

		os.Exit(-1)
	}
}
