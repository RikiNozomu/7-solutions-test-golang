package handler

import (
	domain "7-solutions-test-golang/core/domains"
	service "7-solutions-test-golang/core/services"
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

// Verifies the /auth/login endpoint for both failure and success scenarios.
func TestLoginRoute(t *testing.T) {
	// Set environment variables for JWT configuration
	os.Setenv("JWT_SECRET", "jwt-secret")
	os.Setenv("JWT_TIME_EXPIRED_SECOND", "3600")

	// Setup Gin router with middleware
	router := gin.Default()
	util.RegisterValidation()
	router.Use(middleware.ErrorResponseMiddleware())
	router.Use(middleware.LogRequestMiddleware())

	// Initialize mock user service and auth handler
	userRepo := mock.NewMockUserRepository()
	userService := service.CreateUserService(userRepo)
	authService := service.CreateAuthService(userService)
	authHandler := NewAuthHandler(authService)
	authHandler.AuthRoutes(router)

	// Prepare test credentials
	payload := domain.LoginPayload{
		Email:    "test@gmail.com",
		Password: "psp1234ASAS",
	}

	// Create a user for login testing
	_, err := userService.Create(domain.UserCreate{
		Name:     "tset",
		Email:    payload.Email,
		Password: payload.Password,
	})
	assert.Equal(t, nil, err)

	// Attempt login with incorrect password
	w := httptest.NewRecorder()
	badPayload := map[string]string{
		"email":    payload.Email,
		"password": "random1234",
	}
	jsonBytes, _ := json.Marshal(badPayload)
	req, _ := http.NewRequest("POST", "/auth/login", bytes.NewBuffer(jsonBytes))
	router.ServeHTTP(w, req)
	assert.Equal(t, 401, w.Code)

	// Attempt login with correct credentials
	w = httptest.NewRecorder()
	jsonBytes, _ = json.Marshal(payload)
	req, _ = http.NewRequest("POST", "/auth/login", bytes.NewBuffer(jsonBytes))
	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)

	// Parse and validate response
	var response map[string]any
	json.Unmarshal(w.Body.Bytes(), &response)
	data := response["data"].(map[string]any)
	assert.Greater(t, len(data["token"].(string)), 0)
}
