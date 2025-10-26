package middleware

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// Validates the presence and integrity of a JWT token from the Authorization header.
func CheckBarrierHeader() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract the Authorization header (expected format: "Bearer <token>")
		tokenString := c.GetHeader("Authorization")
		secretKey := []byte(os.Getenv("JWT_SECRET"))

		// Reject request if no token is provided
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"errors": []string{"Missing barrier token"}})
			c.Abort()
			return
		}

		// Remove "Bearer " prefix to isolate the token
		tokenString = tokenString[len("Bearer "):]

		// Parse and validate the JWT token using the secret key
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
			return secretKey, nil
		})

		// Reject if token parsing fails (e.g., malformed, expired, invalid signature)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"errors": []string{err.Error()}})
			c.Abort()
			return
		}

		// Reject if token is not valid
		if !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"errors": []string{"invalid token"}})
			c.Abort()
			return
		}

		// Store token claims in Gin context for downstream access (e.g., user ID, roles)
		c.Set("claims", token.Claims)

		// Continue to next middleware or handler
		c.Next()
	}
}
