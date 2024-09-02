package chat_api

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	amqpConfig "github.com/nihal-ramaswamy/GoChat/internal/amqp"
	"github.com/nihal-ramaswamy/GoChat/internal/constants"
	"github.com/nihal-ramaswamy/GoChat/internal/db"
	"github.com/nihal-ramaswamy/GoChat/internal/dto"
	"github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
)

type SendChatHandler struct {
	dto.HandlerInterface
	log         *zap.Logger
	pdb         *sql.DB
	middlewares []gin.HandlerFunc
	amqpConfig  *amqpConfig.AmqpConfig
}

func NewSendChatHandler(
	pdb *sql.DB,
	log *zap.Logger,
	amqpConfig *amqpConfig.AmqpConfig,
) *SendChatHandler {
	return &SendChatHandler{
		log:         log,
		pdb:         pdb,
		amqpConfig:  amqpConfig,
		middlewares: []gin.HandlerFunc{},
	}
}

func (c *SendChatHandler) Pattern() string {
	return "/chat"
}

func (c *SendChatHandler) RequestMethod() string {
	return constants.POST
}

func (c *SendChatHandler) Middlewares() []gin.HandlerFunc {
	return c.middlewares
}

// Handler to send a chat
// POST /chat/chat
//
//	Request Body: {
//	 "receiverId": receiverId,
//	 "message": message
//	 }
//	 Response:
//	 200 OK: {
//	 "message": "Chat sent"
//	 }
//	 400 Bad Request: {
//	 "error": "Error reading payload"
//	 }
//	 500 Internal Server Error: {
//	 "error": "Error getting user from email"
//	 }
//	 500 Internal Server Error: {
//	 "error": "Error saving chat"
//	 }
//	 500 Internal Server Error: {
//	 "error": "Error marshalling chat"
//	 }
func (c *SendChatHandler) Handler() gin.HandlerFunc {
	return func(ginCtx *gin.Context) {
		email := ginCtx.GetString("email")
		sender, err := db.GetUserFromEmail(c.pdb, email)
		if err != nil {
			c.log.Error("Error getting user from email")
			ginCtx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": "Error getting user from email",
			})
			return
		}

		senderId := sender.Id

		var chat dto.Chat
		if err := ginCtx.ShouldBindJSON(&chat); err != nil {
			c.log.Error("Error binding json", zap.Error(err))
			ginCtx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": "Error reading payload",
			})
			return
		}
		chat.SenderId = senderId
		chat.CreatedAt = time.Now()

		// Save to db
		if err := db.SaveChat(c.pdb, &chat); err != nil {
			c.log.Error("Error saving chat", zap.Error(err))
			ginCtx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": "Error saving chat",
			})
			return
		}

		// Send to rmq queue
		var wg sync.WaitGroup
		wg.Add(1)

		go func() {
			defer wg.Done()
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			body, err := json.Marshal(chat)
			if err != nil {
				c.log.Error("Error marshalling chat", zap.Error(err))
			}

			err = c.amqpConfig.Channel.PublishWithContext(
				ctx,
				constants.EXCHANGE_NAME, // Exchange
				chat.ReceiverId,         // Routing key
				false,                   // mandatory
				false,                   // immediate
				amqp091.Publishing{
					ContentType: "text/plain",
					Body:        body,
				})
			if err != nil {
				c.log.Error("Error publishing message", zap.Error(err))
			}
		}()

		wg.Wait()

		ginCtx.JSON(http.StatusOK, gin.H{
			"message": "Chat sent",
		})
	}
}
