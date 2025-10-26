package handler

import (
	domain "7-solutions-test-golang/core/domains"
	service "7-solutions-test-golang/core/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	service *service.AuthService
}

func NewAuthHandler(service *service.AuthService) *AuthHandler {
	return &AuthHandler{service}
}

func (h *AuthHandler) Login(c *gin.Context) {
	var payload domain.LoginPayload

	if err := c.ShouldBindJSON(&payload); err != nil {
		c.Error(gin.Error{
			Err:  err,
			Type: gin.ErrorTypeBind,
		})
		return
	}
	token, expire_in_unix, err := h.service.Login(payload)
	if err != nil {
		c.Error(gin.Error{
			Err:  err,
			Type: gin.ErrorTypeAny,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": gin.H{"token": token, "expire_in_unix": expire_in_unix}})
}

func (h *AuthHandler) AuthRoutes(router *gin.Engine) {
	authGroup := router.Group("/auth")
	{
		authGroup.POST("/login", h.Login)
	}
}
