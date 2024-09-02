package chat_api

import (
	"database/sql"

	"github.com/gin-gonic/gin"
	"github.com/nihal-ramaswamy/GoChat/internal/constants"
	"github.com/nihal-ramaswamy/GoChat/internal/db"
	"github.com/nihal-ramaswamy/GoChat/internal/dto"
	"go.uber.org/zap"
)

type ReadChatDbHandler struct {
	dto.HandlerInterface
	middleware []gin.HandlerFunc
	pdb        *sql.DB
	log        *zap.Logger
}

func NewReadDbChatHandler(pdb *sql.DB, log *zap.Logger) *ReadChatDbHandler {
	return &ReadChatDbHandler{
		pdb:        pdb,
		log:        log,
		middleware: []gin.HandlerFunc{},
	}
}

func (r *ReadChatDbHandler) Pattern() string {
	return "/read"
}

// Handler reads chat for a user from DB
// GET /chat/read
//
//	Request Header: {
//	  "Token": Bearer token,
//	  }
//
// Response:
//
//	200 OK: {
//	 "messages": [{
//	 "senderId": senderId,
//	 "receiverId": receiverId,
//	 "message": message,
//	 "timestamp": timestamp
//	 }]
//	 }
func (r *ReadChatDbHandler) Handler() gin.HandlerFunc {
	return func(c *gin.Context) {
		email := c.GetString("email")
		user, err := db.GetUserFromEmail(r.pdb, email)
		if err != nil {
			r.log.Error("error getting user", zap.Error(err))
			c.JSON(500, gin.H{"error": "error getting user"})
			return
		}

		id := user.Id

		messages, err := db.ReadChatForUser(r.pdb, id)
		if err != nil {
			r.log.Error("error reading chat", zap.Error(err))
			c.JSON(500, gin.H{"error": "error reading chat"})
		}

		c.JSON(200, messages)
	}
}

func (r *ReadChatDbHandler) RequestMethod() string {
	return constants.GET
}

func (r *ReadChatDbHandler) Middlewares() []gin.HandlerFunc {
	return r.middleware
}
