package middlewares

import (
	"net/http"

	"github.com/antunesgabriel/babytl_backend/configs"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func (c *gin.Context) {
		const BEARER_SCHEMA = "Bearer "

		authorization := c.GetHeader("Authorization")

		if authorization == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "INVALID_AUTHORIZATION",
			})
		}

		token := authorization[len(BEARER_SCHEMA):]

		tk, err := configs.NewJWT().ValidateToken(token)

		if err != nil || !tk.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "INVALID_AUTHORIZATION",
			})
		}

		claim := tk.Claims.(jwt.MapClaims)

		sum := claim["sum"]

		id := uint(sum.(float64))

		c.Set("authId", id)

		c.Next()
	}
}