package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type ErrorMiddleware struct{}

func (er ErrorMiddleware) GetHandlerFunc() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Step1: Process the request first.
		c.Next()

		// Step2: Check if any errors were added to the context
		if len(c.Errors) > 0 {
			// Step3: Use the last error
			err := c.Errors.Last().Err

			// Step4: Respond with a generic error message
			c.JSON(http.StatusInternalServerError, map[string]any{
				"status":  "error",
				"message": err.Error(),
				"data":    nil,
			})
		}

		// Any other steps if no errors are found
	}
}
