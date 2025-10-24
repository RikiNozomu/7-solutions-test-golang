package middleware

import (
	util "7-solutions-test-golang/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func ErrorResponseMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		if len(c.Errors) > 0 {
			// Get the first error
			err := c.Errors[0].Err

			// Determine the appropriate status code and message
			statusCode := http.StatusInternalServerError
			errorMessages := []string{"An unexpected error occurred"}

			if ginerr, ok := err.(gin.Error); ok {
				if ginerr.Type == gin.ErrorTypeBind {
					statusCode = http.StatusBadRequest
					errorMessages = util.GetValidateErrors(ginerr.Err)
				} else {
					switch ginerr.Err {
					case mongo.ErrNoDocuments:
						errorMessages = []string{"User not found."}
						statusCode = http.StatusNotFound
					case util.ErrorAuthenticated:
						errorMessages = []string{ginerr.Err.Error()}
						statusCode = http.StatusUnauthorized
					case util.ErrorDuplicateKey:
						errorMessages = []string{"Cannot add user with exist email."}
						statusCode = http.StatusConflict
					}
				}
			}

			c.JSON(statusCode, gin.H{"errors": errorMessages})
			c.Abort()
		}
	}

}
