package handler

import (
	middleware "7-solutions-test-golang/middlewares"
	util "7-solutions-test-golang/utils"

	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// Verifies that the root route ("/") responds with status 200 and expected payload.
func TestHealthCheckRoute(t *testing.T) {
	// Setup Gin router with middleware and validation
	router := gin.Default()
	util.RegisterValidation()
	router.Use(middleware.ErrorResponseMiddleware())
	router.Use(middleware.LogRequestMiddleware())

	// Register root and fallback routes
	indexHandler := NewIndexHandler()
	indexHandler.IndexHandler(router)

	// Simulate GET request to "/"
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)
	router.ServeHTTP(w, req)

	// Assert HTTP status code
	assert.Equal(t, 200, w.Code)

	// Parse and assert response body
	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "ok", response["status"])
}
