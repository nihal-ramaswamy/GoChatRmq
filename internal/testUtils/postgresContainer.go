package testUtils

import (
	"context"
	"database/sql"
	"fmt"
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

func SetUpPostgresForTesting(ctx context.Context, rootDir string) (*postgres.PostgresContainer, *sql.DB, error) {
	testConfig := GetDbConfig()

	container, err := GetPostgresContainer(testConfig, rootDir, ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get postgres container: %s", err)
	}

	dbURL, err := container.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get connection string: %s", err)
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to open connection: %s", err)
	}

	err = db.Ping()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to ping db: %s", err)
	}

	return container, db, nil
}

func ReadFromUserDb(db *sql.DB, email string) (int, error) {
	cnt := 0
	query := `SELECT COUNT(*) FROM "USER" WHERE EMAIL = $1`
	err := db.QueryRow(query, email).Scan(&cnt)
	if err != nil {
		return -1, err
	}
	return cnt, err
}
