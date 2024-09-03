package auth_api

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nihal-ramaswamy/GoChat/internal/constants"
	"github.com/nihal-ramaswamy/GoChat/internal/db"
	"github.com/nihal-ramaswamy/GoChat/internal/dto"
	"go.uber.org/zap"
)

type NewUserHandler struct {
	dto.HandlerInterface
	db          *sql.DB
	log         *zap.Logger
	middlewares []gin.HandlerFunc
}

func NewNewUserHandler(db *sql.DB, log *zap.Logger) *NewUserHandler {
	return &NewUserHandler{
		db:  db,
		log: log,
	}
}

func (*NewUserHandler) Pattern() string {
	return "/register"
}

// Handler creates a new user in the database
// POST /auth/register
//
//	Request Body: {
//	  "email": email,
//	  "password": password,
//	  "name": name
//	  }
//
// Response:
//
//	202 Accepted: {
//	 "id": id
//	 }
//	 400 Bad Request: {
//	 "error": "User with email %s already exists"
//	 }
func (n *NewUserHandler) Handler() gin.HandlerFunc {
	return func(c *gin.Context) {
		user := dto.NewUser()

		if err := c.ShouldBindJSON(&user); err != nil {
			err := c.Error(err)
			n.log.Info("Responding with error", zap.Error(err))
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		if db.DoesEmailExist(n.db, user.Email) {
			c.JSON(http.StatusBadRequest,
				gin.H{
					"error": fmt.Sprintf("User with email %s already exists", user.Email),
				})
			return
		}

		id := db.RegisterNewUser(n.db, user, n.log)

		c.JSON(http.StatusAccepted, gin.H{"id": id})
	}
}

func (*NewUserHandler) RequestMethod() string {
	return constants.POST
}

func (n *NewUserHandler) Middlewares() []gin.HandlerFunc {
	return n.middlewares
}
