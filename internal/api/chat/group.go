package chat_api

import (
	"context"
	"database/sql"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	amqpConfig "github.com/nihal-ramaswamy/GoChat/internal/amqp"
	"github.com/nihal-ramaswamy/GoChat/internal/dto"
	"github.com/nihal-ramaswamy/GoChat/internal/middlewares"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type ChatGroup struct {
	dto.ServerGroupInterface
	routeHandlers []dto.HandlerInterface
	middlewares   []gin.HandlerFunc
}

func NewChatGroup(
	pdb *sql.DB,
	rdb_auth *redis.Client,
	ctx context.Context,
	log *zap.Logger,
	amqpConfig *amqpConfig.AmqpConfig,
	upgrader *websocket.Upgrader,
	websocketMap *dto.WebsocketConnectionMap,
) *ChatGroup {
	handlers := []dto.HandlerInterface{
		NewSendChatHandler(pdb, log, amqpConfig),
		NewReadDbChatHandler(pdb, log),
		NewReadChatWsHandler(pdb, log, upgrader, websocketMap, amqpConfig),
	}

	return &ChatGroup{
		routeHandlers: handlers,
		middlewares:   []gin.HandlerFunc{middlewares.AuthMiddleware(pdb, rdb_auth, ctx, log)},
	}
}

func (cg *ChatGroup) Group() string {
	return "/chat"
}

func (cg *ChatGroup) RouteHandlers() []dto.HandlerInterface {
	return cg.routeHandlers
}

func (cg *ChatGroup) Middlewares() []gin.HandlerFunc {
	return cg.middlewares
}
