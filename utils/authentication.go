package util

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func IsValid(c *gin.Context) bool {
	claims, exists := c.Get("claims")
	if !exists {
		return false
	}

	userClaims := claims.(jwt.MapClaims)
	id := userClaims["id"].(string)

	return c.Param("id") == id
}
