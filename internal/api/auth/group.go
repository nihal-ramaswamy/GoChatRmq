package auth_api

import (
	"context"
	"database/sql"

	"github.com/gin-gonic/gin"
	"github.com/nihal-ramaswamy/GoChat/internal/dto"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type AuthGroup struct {
	dto.ServerGroupInterface
	routeHandlers []dto.HandlerInterface
	middlewares   []gin.HandlerFunc
}

func (*AuthGroup) Group() string {
	return "/auth"
}

func (h *AuthGroup) RouteHandlers() []dto.HandlerInterface {
	return h.routeHandlers
}

func NewAuthGroup(
	db *sql.DB,
	rdb *redis.Client,
	ctx context.Context,
	log *zap.Logger,
) *AuthGroup {
	handlers := []dto.HandlerInterface{
		NewNewUserHandler(db, log),
		NewLoginUserHandler(db, rdb, ctx, log),
		NewLogoutUserHandler(db, rdb, ctx, log),
	}

	return &AuthGroup{
		routeHandlers: handlers,
		middlewares:   []gin.HandlerFunc{},
	}
}

func (*AuthGroup) AuthRequired() bool {
	return false
}

func (a *AuthGroup) Middlewares() []gin.HandlerFunc {
	return a.middlewares
}
