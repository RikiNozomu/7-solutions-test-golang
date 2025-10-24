package middleware

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

func LogRequestMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		duration := time.Since(start)
		fmt.Printf("%s - %s | Executed Time = %.4fs\n", c.Request.Method, c.Request.URL.Path, duration.Seconds())
	}
}
