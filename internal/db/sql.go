package db

import (
	"database/sql"

	"github.com/nihal-ramaswamy/GoChat/internal/dto"
)

func insertIntoUser(db *sql.DB, user *dto.User) (string, error) {
	if db == nil {
		panic("db cannot be nil")
	}

	var id string
	query := `INSERT INTO "USER" (NAME, EMAIL, PASSWORD) VALUES ($1, $2, $3) RETURNING ID`
	err := db.QueryRow(query, user.Name, user.Email, user.Password).Scan(&id)

	return id, err
}

func selectAllFromUserWhereEmailIs(db *sql.DB, email string) (dto.User, error) {
	if db == nil {
		panic("db cannot be nil")
	}

	var user dto.User
	query := `SELECT ID, NAME, EMAIL FROM "USER" WHERE EMAIL = $1`
	err := db.QueryRow(query, email).Scan(&user.Id, &user.Name, &user.Email)
	if err != nil {
		return user, err
	}

	return user, err
}

func selectPasswordFromUserWhereEmailIDs(db *sql.DB, email string) (string, error) {
	if db == nil {
		panic("db cannot be nil")
	}
	var password string
	query := `SELECT PASSWORD FROM "USER" WHERE EMAIL = $1`
	err := db.QueryRow(query, email).Scan(&password)

	return password, err
}

func insertIntoChat(db *sql.DB, chat *dto.Chat) error {
	if db == nil {
		panic("db cannot be nil")
	}
	query := `INSERT INTO "CHAT" (SENDER_ID, RECEIVER_ID, MESSAGE, CREATED_AT) VALUES ($1, $2, $3, $4)`
	_, err := db.Exec(query, chat.SenderId, chat.ReceiverId, chat.Message, chat.CreatedAt)
	return err
}

func selectAllFromChatWhereUserIdIs(db *sql.DB, id string) ([]dto.Chat, error) {
	if db == nil {
		panic("db cannot be nil")
	}
	var chats []dto.Chat
	query := `SELECT SENDER_ID, RECEIVER_ID, MESSAGE, CREATED_AT FROM "CHAT" WHERE SENDER_ID = $1 OR RECEIVER_ID = $1`
	rows, err := db.Query(query, id)
	defer rows.Close()
	if err != nil {
		return chats, err
	}
	for rows.Next() {
		var chat dto.Chat
		err = rows.Scan(&chat.SenderId, &chat.ReceiverId, &chat.Message, &chat.CreatedAt)
		if err != nil {
			return chats, err
		}
		chats = append(chats, chat)
	}
	return chats, err
}
