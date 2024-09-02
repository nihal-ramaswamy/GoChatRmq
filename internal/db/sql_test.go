package db

import (
	"context"
	"database/sql"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"sync"
	"testing"

	_ "github.com/lib/pq"
	"github.com/nihal-ramaswamy/GoChat/internal/dto"
	"github.com/nihal-ramaswamy/GoChat/internal/utils"
)

func runUserTest(t *testing.T, db *sql.DB) {
	user := &dto.User{
		Name:     utils.RandStringRunes(10),
		Email:    utils.RandStringRunes(10),
		Password: utils.RandStringRunes(10),
	}

	_, err := insertIntoUser(db, user)
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

func runChatTest(t *testing.T, db *sql.DB) {
	userIds := []string{}
	for range 5 {
		userIds = append(userIds, utils.RandStringRunes(10))
	}
	chats := []*dto.Chat{}

	for range 20 {
		userId1 := userIds[rand.Intn(len(userIds))]
		userId2 := userIds[rand.Intn(len(userIds))]
		for userId1 == userId2 {
			userId2 = userIds[rand.Intn(len(userIds))]
		}
		chatDto := &dto.Chat{
			SenderId:   userIds[rand.Intn(len(userIds))],
			ReceiverId: userIds[rand.Intn(len(userIds))],
			Message:    utils.RandStringRunes(10),
		}
		chats = append(chats, chatDto)
	}

	var wg sync.WaitGroup

	for _, chat := range chats {
		wg.Add(1)
		go func(chat *dto.Chat) {
			defer wg.Done()
			err := insertIntoChat(db, chat)
			if err != nil {
				t.Fatalf("Error inserting chat: %s", err)
			}
		}(chat)
	}
	wg.Wait()

	for _, userId := range userIds {
		wg.Add(1)

		go func(userId string) {
			defer wg.Done()
			chatsFromDB, err := selectAllFromChatWhereUserIdIs(db, userId)
			if err != nil {
				t.Fatalf("Error selecting chats: %s", err)
			}

			expectedChats := utils.FilterChatsByUserId(chats, userId)
			if len(chatsFromDB) != expectedChats {
				t.Fatalf("Expected %d chats, got %d", expectedChats, len(chatsFromDB))
			}
		}(userId)
	}

	wg.Wait()
}

// Tests:
// 1. Inserting a user into the database
// 2. Selecting a user name from the database
// 3. Selecting the password of a user
func TestUser(t *testing.T) {
	testConfig := utils.GetDbConfig()
	ctx := context.Background()

	pwd, err := os.Getwd()
	// Go up two directories
	rootDir := filepath.Dir(filepath.Dir(pwd))

	if err != nil {
		t.Fatalf("failed to get working directory: %s", err)
	}

	container, err := utils.GetPostgresContainer(testConfig, rootDir, ctx)
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

	var wg sync.WaitGroup
	for range 10 {
		wg.Add(1)
		go func(t *testing.T) {
			defer wg.Done()
			runUserTest(t, db)
		}(t)
	}
	wg.Wait()
}

// Tests inserting a chat into the database and selecting all chats for each user
func TestChat(t *testing.T) {
	testConfig := utils.GetDbConfig()
	ctx := context.Background()

	pwd, err := os.Getwd()
	// Go up two directories
	rootDir := filepath.Dir(filepath.Dir(pwd))

	if err != nil {
		t.Fatalf("failed to get working directory: %s", err)
	}

	container, err := utils.GetPostgresContainer(testConfig, rootDir, ctx)
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

	runChatTest(t, db)
}
