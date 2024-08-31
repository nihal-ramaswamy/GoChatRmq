package dto

import "github.com/gin-gonic/gin"

type ServerGroupInterface interface {
	Group() string
	RouteHandlers() []HandlerInterface
	Middlewares() []gin.HandlerFunc
}
