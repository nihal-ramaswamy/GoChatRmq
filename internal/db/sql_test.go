package db_test

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
	query "github.com/nihal-ramaswamy/GoChat/internal/db"
	"github.com/nihal-ramaswamy/GoChat/internal/dto"
	"github.com/nihal-ramaswamy/GoChat/internal/testUtils"
	"go.uber.org/zap"
)

func runUserTest(t *testing.T, db *sql.DB) {
	user := &dto.User{
		Name:     testUtils.RandStringRunes(10),
		Email:    testUtils.RandStringRunes(10),
		Password: testUtils.RandStringRunes(10),
	}

	log, err := zap.NewDevelopment()
	if err != nil {
		t.Fatalf("Error creating zap logger: %s", err)
	}
	_ = query.RegisterNewUser(db, user, log)

	userFromDB, err := query.GetUserFromEmail(db, user.Email)
	if err != nil {
		t.Fatalf("Error selecting user: %s", err)
	}
	if userFromDB.Name != user.Name {
		t.Fatalf("Expected %s, got %s", user.Name, userFromDB.Name)
	}

	pass := query.DoesPasswordMatch(db, user, log)
	if pass {
		t.Fatalf("Error matching password. Expected no match")
	}
}

func runChatTest(t *testing.T, db *sql.DB) {
	userIds := []string{}
	for range 5 {
		userIds = append(userIds, testUtils.RandStringRunes(10))
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
			Message:    testUtils.RandStringRunes(10),
		}
		chats = append(chats, chatDto)
	}

	var wg sync.WaitGroup

	for _, chat := range chats {
		wg.Add(1)
		go func(chat *dto.Chat) {
			defer wg.Done()
			err := query.SaveChat(db, chat)
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
			chatsFromDB, err := query.ReadChatForUser(db, userId)
			if err != nil {
				t.Fatalf("Error selecting chats: %s", err)
			}

			expectedChats := testUtils.FilterChatsByUserId(chats, userId)
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
	ctx := context.Background()
	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Error getting working directory: %s", err)
	}
	rootDir := filepath.Join(wd, "..", "..")
	fmt.Println(rootDir)

	container, db, err := testUtils.SetUpPostgresForTesting(ctx, rootDir)
	if err != nil {
		t.Fatalf("Error setting up postgres for testing: %s", err)
	}

	t.Cleanup(func() {
		if err := container.Terminate(ctx); err != nil {
			t.Fatalf("failed to terminate container: %s", err)
		}
		db.Close()
	})

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
	ctx := context.Background()
	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Error getting working directory: %s", err)
	}
	rootDir := filepath.Join(wd, "..", "..")
	fmt.Println(rootDir)

	container, db, err := testUtils.SetUpPostgresForTesting(ctx, rootDir)
	if err != nil {
		t.Fatalf("Error setting up postgres for testing: %s", err)
	}

	t.Cleanup(func() {
		if err := container.Terminate(ctx); err != nil {
			t.Fatalf("failed to terminate container: %s", err)
		}
		db.Close()
	})

	runChatTest(t, db)
}

// Tests password encryption and decryption
// func TestPassword(t *testing.T) {
// 	userDto := &dto.User{
// 		Name:     testUtils.RandStringRunes(10),
// 		Email:    testUtils.RandStringRunes(10),
// 		Password: testUtils.RandStringRunes(10),
// 	}
// 	ctx := context.Background()
// 	wd, err := os.Getwd()
// 	if err != nil {
// 		t.Fatalf("Error getting working directory: %s", err)
// 	}
//
// 	rootDir := filepath.Join(wd, "..", "..")
// 	container, db, err := testUtils.SetUpPostgresForTesting(ctx, rootDir)
// 	if err != nil {
// 		t.Fatalf("Error setting up postgres for testing: %s", err)
// 	}
//
// 	t.Cleanup(func() {
// 		if err := container.Terminate(ctx); err != nil {
// 			t.Fatalf("failed to terminate container: %s", err)
// 		}
// 		db.Close()
// 	})
//
// 	log, err := zap.NewDevelopment()
// 	if err != nil {
// 		t.Fatalf("Error creating zap logger: %s", err)
// 	}
// 	RegisterNewUser(db, userDto, log)
//
// 	if !DoesPasswordMatch(db, userDto, log) {
// 		t.Fatalf("Error matching password. Expected match")
// 	}
//
// 	userDto.Password = testUtils.RandStringRunes(11)
// 	if DoesPasswordMatch(db, userDto, log) {
// 		t.Fatalf("Error matching password. Expected no match")
// 	}
// }
