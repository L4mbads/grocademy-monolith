package middlewares

import (
	"grocademy/internal/auth"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type AuthAPIMiddleware struct{}

func NewAuthAPIMiddleware() *AuthAPIMiddleware {
	return &AuthAPIMiddleware{}
}

func (am AuthAPIMiddleware) GetHandlerFunc() gin.HandlerFunc {
	return func(c *gin.Context) {
		var tokenString string
		var err error

		tokenString, err = c.Cookie("jwt_token")
		if err != nil {
			// if no cookie, eg. from admin FE
			authHeader := c.GetHeader("Authorization")
			if authHeader == "" {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
				return
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid Authorization header format"})
				return
			}
			tokenString = parts[1]
		}

		claims, err := auth.ValidateJWT(tokenString)
		if err != nil {
			status := http.StatusUnauthorized
			if err.Error() == "token expired" {
				status = http.StatusForbidden // Use 403 for expired token, 401 for invalid format/signature
			}
			c.AbortWithStatusJSON(status, gin.H{"error": "Invalid or expired token: " + err.Error()})
			return
		}

		// Store user information in context for handlers
		c.Set("username", claims.Username)
		c.Set("email", claims.Email)
		c.Set("id", claims.ID)
		c.Next()
	}
}
