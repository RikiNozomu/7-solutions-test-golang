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

func TestHealthCheckRoute(t *testing.T) {
	router := gin.Default()
	util.RegisterValidation()
	router.Use(middleware.ErrorResponseMiddleware())
	router.Use(middleware.LogRequestMiddleware())

	indexHandler := NewIndexHandler()
	indexHandler.IndexHandler(router)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	assert.Equal(t, "ok", response["status"])
}
