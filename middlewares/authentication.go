package middleware

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func CheckBarrierHeader() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		secretKey := []byte(os.Getenv("JWT_SECRET"))
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"errors": []string{"Missing barrier token"}})
			c.Abort()
			return
		}

		tokenString = tokenString[len("Bearer "):]
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
			return secretKey, nil
		})

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"errors": []string{err.Error()}})
			c.Abort()
			return
		}

		if !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"errors": []string{"invalid token"}})
			c.Abort()
			return
		}

		c.Set("claims", token.Claims)
		c.Next()
	}
}
