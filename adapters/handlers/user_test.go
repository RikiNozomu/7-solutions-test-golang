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

var name = "testpongpong"
var email = "email@com.com"
var password = "test1234ASAS"

var badEmail = "email"
var badPassword = "tes234"

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

func TestUserGetRoute(t *testing.T) {
	var response map[string]any
	router, userService := setUpServicesAndRoute()
	// Create user before.
	user, _ := userService.Create(domain.UserCreate{
		Name:     name,
		Email:    email,
		Password: password,
	})

	// Get All
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/user", nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
	json.Unmarshal(w.Body.Bytes(), &response)
	datas := response["data"].([]any)
	assert.Greater(t, len(datas), 0)

	// Get not exist
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/user/barbarbar", nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, 404, w.Code)

	// Get exist
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/user/"+user.ID.Hex(), nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
	json.Unmarshal(w.Body.Bytes(), &response)
	data := response["data"].(map[string]any)
	assert.Equal(t, data["email"].(string), email)
}

func TestUserCreateRoute(t *testing.T) {
	router, _ := setUpServicesAndRoute()

	// Incorrect Email payload
	w := httptest.NewRecorder()
	payload := map[string]string{"email": badEmail, "password": badPassword, "name": name}
	jsonBytes, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", "/user", bytes.NewBuffer(jsonBytes))
	router.ServeHTTP(w, req)
	assert.Equal(t, 400, w.Code)

	// Incorrect Email payload
	w = httptest.NewRecorder()
	payload = map[string]string{"email": email, "password": badPassword, "name": name}
	jsonBytes, _ = json.Marshal(payload)
	req, _ = http.NewRequest("POST", "/user", bytes.NewBuffer(jsonBytes))
	router.ServeHTTP(w, req)
	assert.Equal(t, 400, w.Code)

	// Correct Account
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

func TestUserUpdateRoute(t *testing.T) {
	var response map[string]any
	newName := name + "123456"
	router, userService := setUpServicesAndRoute()
	// Create user before.
	user, _ := userService.Create(domain.UserCreate{
		Name:     name,
		Email:    email,
		Password: password,
	})

	// Update without Access token
	w := httptest.NewRecorder()
	payload := map[string]string{"name": newName}
	jsonBytes, _ := json.Marshal(payload)
	req, _ := http.NewRequest("PUT", "/user/"+user.ID.Hex(), bytes.NewBuffer(jsonBytes))
	router.ServeHTTP(w, req)
	assert.Equal(t, 401, w.Code)

	// Login to get access token
	w = httptest.NewRecorder()
	payloadLogin := map[string]string{"email": email, "password": password}
	jsonBytes, _ = json.Marshal(payloadLogin)
	req, _ = http.NewRequest("POST", "/auth/login", bytes.NewBuffer(jsonBytes))
	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
	json.Unmarshal(w.Body.Bytes(), &response)
	data := response["data"].(map[string]any)

	// Update with Access Token
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

func TestUseDeleteRoute(t *testing.T) {
	var response map[string]any
	router, userService := setUpServicesAndRoute()
	// Create user before.
	user, _ := userService.Create(domain.UserCreate{
		Name:     name,
		Email:    email,
		Password: password,
	})

	// Delete without Access token
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/user/"+user.ID.Hex(), nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, 401, w.Code)

	// Login to get access token
	w = httptest.NewRecorder()
	payloadLogin := map[string]string{"email": email, "password": password}
	jsonBytes, _ := json.Marshal(payloadLogin)
	req, _ = http.NewRequest("POST", "/auth/login", bytes.NewBuffer(jsonBytes))
	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
	json.Unmarshal(w.Body.Bytes(), &response)
	data := response["data"].(map[string]any)

	// Delete with Access Token
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("DELETE", "/user/"+user.ID.Hex(), nil)
	req.Header.Set("Authorization", "Bearer "+data["token"].(string))
	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
}
