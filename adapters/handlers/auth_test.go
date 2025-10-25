package handler

import (
	domain "7-solutions-test-golang/core/domains"
	service "7-solutions-test-golang/core/service"
	middleware "7-solutions-test-golang/middlewares"
	mock "7-solutions-test-golang/mocks"
	util "7-solutions-test-golang/utils"
	"bytes"
	"os"

	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestLoginRoute(t *testing.T) {
	os.Setenv("JWT_SECRET", "jwt-secret")
	os.Setenv("JWT_TIME_EXPIRED_SECOND", "3600")

	router := gin.Default()
	util.RegisterValidation()
	router.Use(middleware.ErrorResponseMiddleware())
	router.Use(middleware.LogRequestMiddleware())

	userRepo := mock.NewMockUserRepository()
	userService := service.CreateUserService(userRepo)

	authService := service.CreateAuthService(userService)
	authHandler := NewAuthHandler(authService)
	authHandler.AuthRoutes(router)

	payload := domain.LoginPayload{Email: "test@gmail.com", Password: "psp1234ASAS"}
	_, err := userService.Create(domain.UserCreate{
		Name:     "tset",
		Email:    payload.Email,
		Password: payload.Password,
	})
	assert.Equal(t, nil, err)

	// Incorrect Account
	w := httptest.NewRecorder()
	badPayload := map[string]string{"email": payload.Email, "password": "random1234"}
	jsonBytes, _ := json.Marshal(badPayload)
	req, _ := http.NewRequest("POST", "/auth/login", bytes.NewBuffer(jsonBytes))
	router.ServeHTTP(w, req)
	assert.Equal(t, 401, w.Code)

	// Correct Account
	w = httptest.NewRecorder()
	jsonBytes, _ = json.Marshal(payload)
	req, _ = http.NewRequest("POST", "/auth/login", bytes.NewBuffer(jsonBytes))
	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)

	var response map[string]any
	json.Unmarshal(w.Body.Bytes(), &response)
	data := response["data"].(map[string]any)
	assert.Greater(t, len(data["token"].(string)), 0)
}
