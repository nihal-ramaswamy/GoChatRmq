package utils

import (
	"context"
	"math/rand"
	"path/filepath"
	"time"

	"github.com/nihal-ramaswamy/GoChat/internal/dto"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

func GetDbConfig() *dto.TestConfigDto {
	return &dto.TestConfigDto{
		Username:     "postgresTest",
		Password:     "postgresTest",
		DatabaseName: "go_chat",
	}
}

func GetPostgresContainer(
	testConfig *dto.TestConfigDto,
	rootDir string,
	ctx context.Context,
) (*postgres.PostgresContainer, error) {
	container, err := postgres.Run(
		ctx,
		"docker.io/postgres:16-alpine",
		postgres.WithInitScripts(filepath.Join(rootDir, "db", "init.sql")),
		postgres.WithUsername(testConfig.Username),
		postgres.WithPassword(testConfig.Password),
		postgres.WithDatabase(testConfig.DatabaseName),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).WithStartupTimeout(5*time.Second)),
		// postgres.WithSQLDriver("pq"),
	)

	return container, err
}

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
