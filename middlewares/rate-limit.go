package middleware

import (
	"net/http"
	"os"
	"strconv"

	"golang.org/x/time/rate"

	"github.com/gin-gonic/gin"
)

func RateLimiter() gin.HandlerFunc {
	limit, err := strconv.Atoi(os.Getenv("LIMIT_RATE_IN_SECOND"))
	if err != nil {
		limit = 10
	}
	if limit <= 0 {
		limit = 10
	}

	brust, err := strconv.Atoi(os.Getenv("MAXIMUM_BRUST"))
	if err != nil {
		brust = 100
	}
	if brust <= 0 {
		brust = 100
	}
	limiter := rate.NewLimiter(rate.Limit(limit), brust)
	return func(c *gin.Context) {

		if limiter.Allow() {
			c.Next()
		} else {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": []string{"Limite exceed"},
			})
		}

	}
}
