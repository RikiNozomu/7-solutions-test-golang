package handler

import (
	domain "7-solutions-test-golang/core/domains"
	"7-solutions-test-golang/core/service"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type AuthHandler struct {
	service *service.AuthService
}

func NewAuthHandler(service *service.AuthService) *AuthHandler {
	return &AuthHandler{service}
}

func (h *AuthHandler) Login(c *gin.Context) {
	var payload domain.LoginPayload
	secretKey := []byte(os.Getenv("JWT_SECRET"))
	expSecond, err := strconv.Atoi(os.Getenv("JWT_TIME_EXPIRED_SECOND"))
	if err != nil {
		c.Error(gin.Error{
			Err:  err,
			Type: gin.ErrorTypeAny,
		})
	}

	if err := c.ShouldBindJSON(&payload); err != nil {
		c.Error(gin.Error{
			Err:  err,
			Type: gin.ErrorTypeBind,
		})
		return
	}
	user, err := h.service.Login(payload)
	if err != nil {
		c.Error(gin.Error{
			Err:  err,
			Type: gin.ErrorTypeAny,
		})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"id":  user.ID.Hex(),
			"exp": time.Now().Add(time.Duration(expSecond) * time.Second).Unix(),
		})

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		c.Error(gin.Error{
			Err:  err,
			Type: gin.ErrorTypeAny,
		})
	}

	c.JSON(http.StatusOK, gin.H{"token": tokenString, "expire_in_unix": time.Now().Add(time.Duration(expSecond) * time.Second).Unix()})
}

func (h *AuthHandler) AuthRoutes(router *gin.Engine) {
	authGroup := router.Group("/auth")
	{
		authGroup.POST("/login", h.Login)
	}
}
