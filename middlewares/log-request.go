package middleware

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

// Logs the HTTP method, request path, and execution time for each request.
func LogRequestMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Record the start time before processing the request
		start := time.Now()

		// Proceed to the next middleware or handler
		c.Next()

		// Calculate the duration after the request is processed
		duration := time.Since(start)

		// Log the method, path, and execution time in seconds
		fmt.Printf("%s - %s | Executed Time = %.4fs\n", c.Request.Method, c.Request.URL.Path, duration.Seconds())
	}
}
