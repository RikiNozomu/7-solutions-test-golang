package handler

import (
	domain "7-solutions-test-golang/core/domains"
	service "7-solutions-test-golang/core/services"
	middleware "7-solutions-test-golang/middlewares"
	mock "7-solutions-test-golang/mocks"
	util "7-solutions-test-golang/utils"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// Mock the datas.
var name = "testpongpong"
var email = "email@com.com"
var password = "test1234ASAS"
var badEmail = "email"
var badPassword = "tes234"

// Setup helper for initializing router and services
func setUpServicesAndRoute() (*gin.Engine, *service.UserService) {
	os.Setenv("JWT_SECRET", "jwt-secret")
	os.Setenv("JWT_TIME_EXPIRED_SECOND", "3600")

	router := gin.Default()
	util.RegisterValidation()
	router.Use(middleware.ErrorResponseMiddleware())
	router.Use(middleware.LogRequestMiddleware())

	userRepo := mock.NewMockUserRepository()
	userService := service.CreateUserService(userRepo)
	userHandler := NewUserHandler(userService)
	userHandler.UserRoutes(router)

	authService := service.CreateAuthService(userService)
	authHandler := NewAuthHandler(authService)
	authHandler.AuthRoutes(router)

	return router, userService
}

// Test GET /user and GET /user/:id
func TestUserGetRoute(t *testing.T) {
	var response map[string]any
	router, userService := setUpServicesAndRoute()

	// Create user
	user, _ := userService.Create(domain.UserCreate{
		Name:     name,
		Email:    email,
		Password: password,
	})

	// GET /user (list all)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/user", nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
	json.Unmarshal(w.Body.Bytes(), &response)
	datas := response["data"].([]any)
	assert.Greater(t, len(datas), 0)

	// GET /user/:id (not found)
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/user/barbarbar", nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, 404, w.Code)

	// GET /user/:id (found)
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/user/"+user.ID.Hex(), nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
	json.Unmarshal(w.Body.Bytes(), &response)
	data := response["data"].(map[string]any)
	assert.Equal(t, data["email"].(string), email)
}

// Test POST /user
func TestUserCreateRoute(t *testing.T) {
	router, _ := setUpServicesAndRoute()

	// Invalid email format
	w := httptest.NewRecorder()
	payload := map[string]string{"email": badEmail, "password": badPassword, "name": name}
	jsonBytes, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", "/user", bytes.NewBuffer(jsonBytes))
	router.ServeHTTP(w, req)
	assert.Equal(t, 400, w.Code)

	// Invalid password format
	w = httptest.NewRecorder()
	payload = map[string]string{"email": email, "password": badPassword, "name": name}
	jsonBytes, _ = json.Marshal(payload)
	req, _ = http.NewRequest("POST", "/user", bytes.NewBuffer(jsonBytes))
	router.ServeHTTP(w, req)
	assert.Equal(t, 400, w.Code)

	// Valid payload
	w = httptest.NewRecorder()
	payload = map[string]string{"email": email, "password": password, "name": name}
	jsonBytes, _ = json.Marshal(payload)
	req, _ = http.NewRequest("POST", "/user", bytes.NewBuffer(jsonBytes))
	router.ServeHTTP(w, req)
	assert.Equal(t, 201, w.Code)

	var response map[string]any
	json.Unmarshal(w.Body.Bytes(), &response)
	data := response["data"].(map[string]any)
	assert.Greater(t, len(data["id"].(string)), 0)
}

// Test PUT /user/:id
func TestUserUpdateRoute(t *testing.T) {
	var response map[string]any
	newName := name + "123456"
	router, userService := setUpServicesAndRoute()

	// Create user
	user, _ := userService.Create(domain.UserCreate{
		Name:     name,
		Email:    email,
		Password: password,
	})

	// Update without token
	w := httptest.NewRecorder()
	payload := map[string]string{"name": newName}
	jsonBytes, _ := json.Marshal(payload)
	req, _ := http.NewRequest("PUT", "/user/"+user.ID.Hex(), bytes.NewBuffer(jsonBytes))
	router.ServeHTTP(w, req)
	assert.Equal(t, 401, w.Code)

	// Login to get token
	w = httptest.NewRecorder()
	payloadLogin := map[string]string{"email": email, "password": password}
	jsonBytes, _ = json.Marshal(payloadLogin)
	req, _ = http.NewRequest("POST", "/auth/login", bytes.NewBuffer(jsonBytes))
	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
	json.Unmarshal(w.Body.Bytes(), &response)
	data := response["data"].(map[string]any)

	// Update with token
	w = httptest.NewRecorder()
	jsonBytes, _ = json.Marshal(payload)
	req, _ = http.NewRequest("PUT", "/user/"+user.ID.Hex(), bytes.NewBuffer(jsonBytes))
	req.Header.Set("Authorization", "Bearer "+data["token"].(string))
	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
	json.Unmarshal(w.Body.Bytes(), &response)
	data = response["data"].(map[string]any)
	assert.Equal(t, newName, data["name"].(string))
}

// Test DELETE /user/:id
func TestUseDeleteRoute(t *testing.T) {
	var response map[string]any
	router, userService := setUpServicesAndRoute()

	// Create user
	user, _ := userService.Create(domain.UserCreate{
		Name:     name,
		Email:    email,
		Password: password,
	})

	// Delete without token
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/user/"+user.ID.Hex(), nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, 401, w.Code)

	// Login to get token
	w = httptest.NewRecorder()
	payloadLogin := map[string]string{"email": email, "password": password}
	jsonBytes, _ := json.Marshal(payloadLogin)
	req, _ = http.NewRequest("POST", "/auth/login", bytes.NewBuffer(jsonBytes))
	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
	json.Unmarshal(w.Body.Bytes(), &response)
	data := response["data"].(map[string]any)

	// Delete with token
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("DELETE", "/user/"+user.ID.Hex(), nil)
	req.Header.Set("Authorization", "Bearer "+data["token"].(string))
	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
}
