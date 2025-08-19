package api

import (
	"grocademy/internal/api/handlers"

	"github.com/gin-gonic/gin"
)

func SetupRouter(userHandler *handlers.UserHandler) *gin.Engine {
	r := gin.Default()

	apiV1 := r.Group("/api")
	{
		users := apiV1.Group("/users")
		{
			users.POST("/", userHandler.CreateUser)
			users.GET("/:id", userHandler.GetUserByID)
		}
	}

	return r
}
