package dto

import "time"

type Chat struct {
	SenderId   string    `json:"sender_id"`
	ReceiverId string    `json:"receiver_id"`
	Message    string    `json:"message"`
	CreatedAt  time.Time `json:"created_at"`
}
