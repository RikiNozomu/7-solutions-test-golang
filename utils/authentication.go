package util

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// IsValid checks whether the authenticated user's ID matches the ID in the request path.
func IsValid(c *gin.Context) bool {
	// Retrieve JWT claims from Gin context (set by authentication middleware)
	claims, exists := c.Get("claims")
	if !exists {
		return false
	}

	// Type assert claims to jwt.MapClaims
	userClaims := claims.(jwt.MapClaims)

	// Extract user ID from JWT claims
	id := userClaims["id"].(string)

	return c.Param("id") == id
}
