package handler

import (
	"github.com/gin-gonic/gin"
)

type IndexHandler struct{}

// NewIndexHandler initializes a new IndexHandler instance.
func NewIndexHandler() *IndexHandler {
	return &IndexHandler{}
}

// IndexHandler registers the root ("/") and fallback (404) routes.
func (h *IndexHandler) IndexHandler(router *gin.Engine) {
	// GET / — Health check or welcome endpoint
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})

	// Catch-all for undefined routes — returns 404 Not Found
	router.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{
			"error": "Not Found",
		})
	})
}
