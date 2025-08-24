package middlewares

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AdminMiddleware struct{}

func NewAdminMiddleware() *AdminMiddleware {
	return &AdminMiddleware{}
}

func (am AdminMiddleware) GetHandlerFunc() gin.HandlerFunc {
	return func(c *gin.Context) {
		username, _ := c.Get("username")
		if username != "admin" {
			c.AbortWithError(http.StatusUnauthorized, errors.New("route only for admin"))
			return
		}
		c.Next()
	}
}
