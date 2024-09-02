package auth_api

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nihal-ramaswamy/GoChat/internal/constants"
	"github.com/nihal-ramaswamy/GoChat/internal/db"
	"github.com/nihal-ramaswamy/GoChat/internal/dto"
	"github.com/nihal-ramaswamy/GoChat/internal/utils"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type LoginUserHandler struct {
	dto.HandlerInterface
	log         *zap.Logger
	db          *sql.DB
	rdb         *redis.Client
	ctx         context.Context
	middlewares []gin.HandlerFunc
}

func NewLoginUserHandler(
	db *sql.DB,
	rdb *redis.Client,
	ctx context.Context,
	log *zap.Logger,
) *LoginUserHandler {
	return &LoginUserHandler{
		log:         log,
		db:          db,
		rdb:         rdb,
		ctx:         ctx,
		middlewares: []gin.HandlerFunc{},
	}
}

func (*LoginUserHandler) Pattern() string {
	return "/signin"
}

func (*LoginUserHandler) RequestMethod() string {
	return constants.POST
}

// Hanlder to authenticate a user
// POST /auth/signin
//
//	Request Body: {
//	  "email": email,
//	  "password": password
//	  }
//
//	  Response:
//	  202 Accepted: {
//	  "token": token
//	  }
//	  401 Unauthorized: {
//	  "error": "Invalid credentials"
//	  }
//	  401 Unauthorized: {
//	  "error": "User with email %s does not exist"
//	  }
//	  500 Internal Server Error: {
//	  "error": "Internal Server Error"
//	  }
//	  400 Bad Request: {
//	  }
func (l *LoginUserHandler) Handler() gin.HandlerFunc {
	return func(c *gin.Context) {
		user := dto.NewUser()
		if err := c.ShouldBindJSON(&user); nil != err {
			err := c.Error(err)
			l.log.Info("Responding with error", zap.Error(err))
			c.AbortWithStatus(http.StatusBadRequest)
		}
		if !db.DoesEmailExist(l.db, user.Email) {
			c.JSON(http.StatusUnauthorized,
				gin.H{
					"error": fmt.Sprintf("User with email %s does not exist", user.Email),
				})
			return
		}

		if !db.DoesPasswordMatch(l.db, user, l.log) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
			return
		}

		token, err := utils.GenerateToken(user)
		if nil != err {
			err := c.Error(err)
			l.log.Info("Responding with error", zap.Error(err))

			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		l.rdb.Set(l.ctx, user.Email, token, constants.TOKEN_EXPIRY_TIME)

		c.JSON(http.StatusAccepted, gin.H{"token": token})
	}
}

func (l *LoginUserHandler) Middlewares() []gin.HandlerFunc {
	return l.middlewares
}
