package middlewares

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/nihal-ramaswamy/GoChat/internal/constants"
	"github.com/nihal-ramaswamy/GoChat/internal/utils"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

func AuthMiddleware(
	pdb *sql.DB,
	rdb *redis.Client,
	ctx context.Context,
	log *zap.Logger,
) gin.HandlerFunc {
	secret := utils.GetDotEnvVariable("SECRET_KEY")
	signingKey := []byte(secret)

	return func(c *gin.Context) {
		token := c.GetHeader("token")
		if token == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized,
				gin.H{"error": "Invalid token. Token should be 'Bearer <Token>'"})
			return
		}
		splitToken := strings.Split(token, constants.BEARER)
		token = splitToken[1]

		if len(splitToken) != 2 {
			c.AbortWithStatusJSON(http.StatusUnauthorized,
				gin.H{"error": "Invalid token. Token should be 'Bearer <Token>'"})
			return
		}

		parsedToken, err := jwt.Parse(
			token,
			func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("Error in parsing token")
				}
				return signingKey, nil
			})

		if nil != err {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		email := parsedToken.Claims.(jwt.MapClaims)["email"].(string)

		_, err = rdb.Get(ctx, email).Result()
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		if claims, ok := parsedToken.Claims.(jwt.MapClaims); ok && parsedToken.Valid {
			c.Set("email", claims["email"])
			c.Set("authenticated", true)
		}

		c.Next()
	}
}
