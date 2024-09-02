package healthcheck_api

import (
	"context"
	"database/sql"

	"github.com/gin-gonic/gin"
	"github.com/nihal-ramaswamy/GoChat/internal/dto"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type HealthCheckGroup struct {
	dto.ServerGroupInterface
	routeHandlers []dto.HandlerInterface
	middlewares   []gin.HandlerFunc
}

func (*HealthCheckGroup) Group() string {
	return "/healthcheck"
}

func (h *HealthCheckGroup) RouteHandlers() []dto.HandlerInterface {
	return h.routeHandlers
}

func NewHealthCheckGroup(
	pdb *sql.DB,
	rdb *redis.Client,
	ctx context.Context,
	log *zap.Logger,
) *HealthCheckGroup {
	handlers := []dto.HandlerInterface{
		NewHealthCheckHandler(),
		NewHealthCheckHandlerAuth(pdb, rdb, ctx, log),
	}

	return &HealthCheckGroup{
		routeHandlers: handlers,
		middlewares:   []gin.HandlerFunc{},
	}
}

func (*HealthCheckGroup) AuthRequired() bool {
	return false
}

func (h *HealthCheckGroup) Middlewares() []gin.HandlerFunc {
	return h.middlewares
}
