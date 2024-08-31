package healthcheck_api

import (
	"github.com/gin-gonic/gin"
	"github.com/nihal-ramaswamy/GoChat/internal/dto"
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

func NewHealthCheckGroup() *HealthCheckGroup {
	handlers := []dto.HandlerInterface{
		NewHealthCheckHandler(),
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
