package healthcheck_api

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

type HealthCheckHandlerAuth struct {
	dto.HandlerInterface
	middlewares []gin.HandlerFunc
}

func NewHealthCheckHandlerAuth(
	pdb *sql.DB,
	rdb *redis.Client,
	ctx context.Context,
	log *zap.Logger,
) *HealthCheckHandlerAuth {
	return &HealthCheckHandlerAuth{
		middlewares: []gin.HandlerFunc{middlewares.AuthMiddleware(pdb, rdb, ctx, log)},
	}
}

func (*HealthCheckHandlerAuth) Pattern() string {
	return "/healthcheckAuth"
}

// Handler returns a handler function for the healthcheck endpoint
// GET /healthcheck/healthcheck
//
//	Request Header: {
//	 "Token": Bearer token,
//	 }
//
// Response Body:
//
//	200 OK: {
//		"message": "Authenticated"
//		}
func (*HealthCheckHandlerAuth) Handler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "Authenticated",
		})
	}
}

func (*HealthCheckHandlerAuth) RequestMethod() string {
	return constants.GET
}

func (h *HealthCheckHandlerAuth) Middlewares() []gin.HandlerFunc {
	return h.middlewares
}
