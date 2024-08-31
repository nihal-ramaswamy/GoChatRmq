package healthcheck_api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nihal-ramaswamy/GoChat/internal/constants"
	"github.com/nihal-ramaswamy/GoChat/internal/dto"
)

type HealthCheckHandler struct {
	dto.HandlerInterface
	middlewares []gin.HandlerFunc
}

func NewHealthCheckHandler() *HealthCheckHandler {
	return &HealthCheckHandler{
		middlewares: []gin.HandlerFunc{},
	}
}

func (*HealthCheckHandler) Pattern() string {
	return "/healthcheck"
}

func (*HealthCheckHandler) Handler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "ok",
		})
	}
}

func (*HealthCheckHandler) RequestMethod() string {
	return constants.GET
}

func (h *HealthCheckHandler) Middlewares() []gin.HandlerFunc {
	return h.middlewares
}
