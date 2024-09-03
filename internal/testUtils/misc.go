package testUtils

import (
	"math/rand"

	"github.com/nihal-ramaswamy/GoChat/internal/dto"
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
