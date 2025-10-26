package handler

import (
	domain "7-solutions-test-golang/core/domains"
	service "7-solutions-test-golang/core/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Handles authentication-related HTTP requests (e.g., login).
type AuthHandler struct {
	service *service.AuthService // Business logic for authentication
}

// NewAuthHandler initializes a new AuthHandler with the given AuthService.
func NewAuthHandler(service *service.AuthService) *AuthHandler {
	return &AuthHandler{service}
}

// Validates credentials and returns a JWT token with expiration timestamp.
func (h *AuthHandler) Login(c *gin.Context) {
	var payload domain.LoginPayload

	// Bind and validate incoming JSON payload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.Error(gin.Error{
			Err:  err,
			Type: gin.ErrorTypeBind,
		})
		return
	}

	// Authenticate user and generate token
	token, expire_in_unix, err := h.service.Login(payload)
	if err != nil {
		c.Error(gin.Error{
			Err:  err,
			Type: gin.ErrorTypeAny,
		})
		return
	}

	// Return token and expiration in response
	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"token":          token,
			"expire_in_unix": expire_in_unix,
		},
	})
}

// AuthRoutes registers authentication-related routes under /auth
func (h *AuthHandler) AuthRoutes(router *gin.Engine) {
	authGroup := router.Group("/auth")
	{
		authGroup.POST("/login", h.Login)
	}
}
