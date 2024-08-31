package auth_api

import (
	"context"
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nihal-ramaswamy/GoChat/internal/constants"
	"github.com/nihal-ramaswamy/GoChat/internal/dto"
	"github.com/nihal-ramaswamy/GoChat/internal/middlewares"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type LogoutUserHandler struct {
	dto.HandlerInterface
	ctx         context.Context
	rdb         *redis.Client
	log         *zap.Logger
	middlewares []gin.HandlerFunc
}

func NewLogoutUserHandler(
	pdb *sql.DB,
	rdb *redis.Client,
	ctx context.Context,
	log *zap.Logger,
) *LogoutUserHandler {
	return &LogoutUserHandler{
		rdb:         rdb,
		ctx:         ctx,
		log:         log,
		middlewares: []gin.HandlerFunc{middlewares.AuthMiddleware(pdb, rdb, ctx, log)},
	}
}

func (*LogoutUserHandler) Pattern() string {
	return "/signout"
}

func (*LogoutUserHandler) RequestMethod() string {
	return constants.POST
}

func (l *LogoutUserHandler) Handler() gin.HandlerFunc {
	return func(c *gin.Context) {
		email := c.GetString("email")
		_, err := l.rdb.Del(l.ctx, email).Result()
		if err != nil {
			l.log.Error("Error deleting token from rdb")
		}

		c.JSON(http.StatusAccepted, gin.H{"message": "ok"})
	}
}

func (l *LogoutUserHandler) Middlewares() []gin.HandlerFunc {
	return l.middlewares
}
