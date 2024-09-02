package routes

import (
	"context"
	"database/sql"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	amqpConfig "github.com/nihal-ramaswamy/GoChat/internal/amqp"
	auth_api "github.com/nihal-ramaswamy/GoChat/internal/api/auth"
	chat_api "github.com/nihal-ramaswamy/GoChat/internal/api/chat"
	healthcheck_api "github.com/nihal-ramaswamy/GoChat/internal/api/healthcheck"
	"github.com/nihal-ramaswamy/GoChat/internal/constants"
	"github.com/nihal-ramaswamy/GoChat/internal/dto"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

func NewRoutes(
	server *gin.Engine,
	pdb *sql.DB,
	rdb_auth *redis.Client,
	ctx context.Context,
	log *zap.Logger,
	amqpConfig *amqpConfig.AmqpConfig,
	upgrader *websocket.Upgrader,
	websocketMap *dto.WebsocketConnectionMap,
) {
	serverGroupHandlers := []dto.ServerGroupInterface{
		healthcheck_api.NewHealthCheckGroup(pdb, rdb_auth, ctx, log),
		auth_api.NewAuthGroup(pdb, rdb_auth, ctx, log),
		chat_api.NewChatGroup(pdb, rdb_auth, ctx, log, amqpConfig, upgrader, websocketMap),
	}

	for _, serverGroupHandler := range serverGroupHandlers {
		newGroup(server, serverGroupHandler)
	}
}

func newGroup(server *gin.Engine, groupHandler dto.ServerGroupInterface) {
	group := server.Group(groupHandler.Group(), groupHandler.Middlewares()...)
	{
		for _, route := range groupHandler.RouteHandlers() {
			newRoute(group, route)
		}
	}
}

func newRoute(server *gin.RouterGroup, routeHandler dto.HandlerInterface) {
	middlewares := routeHandler.Middlewares()
	middlewares = append(middlewares, routeHandler.Handler())
	switch routeHandler.RequestMethod() {
	case constants.GET:
		server.GET(routeHandler.Pattern(), middlewares...)
	case constants.POST:
		server.POST(routeHandler.Pattern(), middlewares...)
	}
}
