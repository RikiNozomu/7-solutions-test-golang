package middleware

import (
	util "7-solutions-test-golang/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

// Converts internal errors into structured JSON responses with appropriate HTTP status codes.
func ErrorResponseMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Continue to next middleware or handler
		c.Next()

		// If any errors occurred during request processing
		if len(c.Errors) > 0 {
			// Extract the first error from Gin's error list
			err := c.Errors[0].Err

			// Default response: internal server error
			statusCode := http.StatusInternalServerError
			errorMessages := []string{"An unexpected error occurred"}

			// Check if the error is a Gin error with metadata
			if ginerr, ok := err.(gin.Error); ok {
				// Handle binding/validation errors
				if ginerr.Type == gin.ErrorTypeBind {
					statusCode = http.StatusBadRequest
					errorMessages = util.GetValidateErrors(ginerr.Err)
				} else {
					// Handle known application-level errors
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

			// Send structured JSON error response
			c.JSON(statusCode, gin.H{"errors": errorMessages})
			c.Abort() // Stop further processing
		}
	}
}
