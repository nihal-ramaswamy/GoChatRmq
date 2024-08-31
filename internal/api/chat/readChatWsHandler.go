package chat_api

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	amqpConfig "github.com/nihal-ramaswamy/GoChat/internal/amqp"
	"github.com/nihal-ramaswamy/GoChat/internal/constants"
	"github.com/nihal-ramaswamy/GoChat/internal/db"
	"github.com/nihal-ramaswamy/GoChat/internal/dto"
	"github.com/nihal-ramaswamy/GoChat/internal/utils"
	"go.uber.org/zap"
)

type ReadChatWsHandler struct {
	dto.HandlerInterface
	middleware   []gin.HandlerFunc
	pdb          *sql.DB
	log          *zap.Logger
	upgrader     *websocket.Upgrader
	websocketMap *dto.WebsocketConnectionMap
	amqpConfig   *amqpConfig.AmqpConfig
}

func NewReadChatWsHandler(
	pdb *sql.DB,
	log *zap.Logger,
	upgrader *websocket.Upgrader,
	websocketMap *dto.WebsocketConnectionMap,
	amqpConfig *amqpConfig.AmqpConfig,
) *ReadChatWsHandler {
	return &ReadChatWsHandler{
		pdb:          pdb,
		log:          log,
		upgrader:     upgrader,
		websocketMap: websocketMap,
		amqpConfig:   amqpConfig,
	}
}

func (r *ReadChatWsHandler) Pattern() string {
	return "/ws"
}

func (r *ReadChatWsHandler) Handler() gin.HandlerFunc {
	return func(c *gin.Context) {
		email := c.GetString("email")
		user, err := db.GetUserFromEmail(r.pdb, email)
		if err != nil {
			r.log.Error("Error getting user from email")
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Error getting user from email",
			})
			return
		}
		id := user.Id

		if err := utils.CreateNewConnection(r.websocketMap, r.upgrader, c, id); err != nil {
			r.log.Error("Error upgrading to websocket connection")
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Error upgrading to websocket connection",
			})
		}

		conn, _ := r.websocketMap.Get(id)

		go func(conn *dto.WebsocketConnection, c *gin.Context) {
			defer conn.Close()

			err = r.amqpConfig.Channel.QueueBind(
				r.amqpConfig.Queue.Name,
				id,                      // routing key
				constants.EXCHANGE_NAME, // exchange
				false,
				nil)
			if err != nil {
				r.log.Error("Error binding queue", zap.Error(err))
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"error": "Error binding queue",
				})
				return
			}

			msgs, err := r.amqpConfig.Channel.Consume(
				r.amqpConfig.Queue.Name,
				"",    // consumer
				true,  // auto ack
				false, // exclusive
				false, // no local
				false, // no wait
				nil,   // args
			)
			if err != nil {
				r.log.Error("Error consuming messages", zap.Error(err))
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"error": "Error consuming messages",
				})
				return
			}
			var forever chan struct{}
			go func() {
				for d := range msgs {
					conn.Conn.WriteMessage(websocket.TextMessage, d.Body)
				}
			}()
			<-forever
		}(conn, c)
	}
}

func (r *ReadChatWsHandler) RequestMethod() string {
	return constants.GET
}

func (r *ReadChatWsHandler) Middlewares() []gin.HandlerFunc {
	return r.middleware
}
