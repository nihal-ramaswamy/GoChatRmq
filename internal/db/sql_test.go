package db

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	_ "github.com/lib/pq"
	"github.com/nihal-ramaswamy/GoChat/internal/dto"
	"github.com/testcontainers/testcontainers-go"
	postgres "github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	"go.uber.org/zap"
)

func TestUser(t *testing.T) {
	username := "postgresTest"
	password := "postgresTest"

	ctx := context.Background()

	log := zap.Must(zap.NewDevelopment())

	pwd, err := os.Getwd()
	// Go up two directories
	rootDir := filepath.Dir(filepath.Dir(pwd))

	if err != nil {
		log.Fatal("failed to get working directory", zap.Error(err))
	}

	container, err := postgres.Run(
		ctx,
		"docker.io/postgres:16-alpine",
		postgres.WithInitScripts(filepath.Join(rootDir, "db", "init.sql")),
		postgres.WithUsername(username),
		postgres.WithPassword(password),
		postgres.WithDatabase("go_chat"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).WithStartupTimeout(5*time.Second)),
		// postgres.WithSQLDriver("pq"),
	)
	if err != nil {
		t.Fatalf("failed to start container, %s", err)
	}

	// Clean up the container
	defer func() {
		if err := container.Terminate(ctx); err != nil {
			t.Fatalf("failed to terminate container: %s", err)
		}
	}()

	dbURL, err := container.ConnectionString(ctx, "sslmode=disable")
	fmt.Println(dbURL)
	if err != nil {
		t.Fatal(err)
	}

	db, err := sql.Open("postgres", dbURL)
	t.Cleanup(func() { db.Close() })

	if err != nil {
		t.Fatalf("Error opening connection: %s", err)
	}

	err = db.Ping()
	if err != nil {
		t.Fatalf("Error pinging db: %s", err)
	}

	user := &dto.User{
		Name:     "John Doe",
		Email:    "jhon",
		Password: "password",
	}

	_, err = insertIntoUser(db, user)
	if err != nil {
		t.Fatalf("Error inserting user: %s", err)
	}

	userFromDB, err := selectAllFromUserWhereEmailIs(db, user.Email)
	if err != nil {
		t.Fatalf("Error selecting user: %s", err)
	}
	if userFromDB.Name != user.Name {
		t.Fatalf("Expected %s, got %s", user.Name, userFromDB.Name)
	}

	pass, err := selectPasswordFromUserWhereEmailIDs(db, user.Email)
	if err != nil {
		t.Fatalf("Error selecting password: %s", err)
	}
	if pass != user.Password {
		t.Fatalf("Expected %s, got %s", user.Password, pass)
	}
}
