package middlewares

import (
	"net/http"

	"grocademy/internal/auth"

	"github.com/gin-gonic/gin"
)

type AuthWebMiddleware struct{}

func NewAuthWebMiddleware() *AuthWebMiddleware {
	return &AuthWebMiddleware{}
}

func (am AuthWebMiddleware) GetHandlerFunc() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString, err := c.Cookie("jwt_token")
		if err != nil {
			// If cookie is not found, redirect to the login page
			c.Redirect(http.StatusFound, "/login")
			c.Abort()
			return
		}

		_, err = auth.ValidateJWT(tokenString)
		if err != nil {
			// If token is invalid or expired, redirect to the login page
			c.Redirect(http.StatusFound, "/login")
			c.Abort()
			return
		}

		c.Next()
	}
}
