package middlewares

import (
	"github.com/gin-gonic/gin"
)

type GinMiddleware interface {
	GetHandlerFunc() gin.HandlerFunc
}
