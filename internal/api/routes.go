package api

import (
	"net/http"

	"grocademy/internal/api/handlers"
	"grocademy/internal/api/middlewares"

	"github.com/gin-gonic/gin"
)

// SetupRouter configures the Gin router with API routes.
func SetupRouter(userHandler *handlers.UserHandler, authHandler *handlers.AuthHandler) *gin.Engine {
	r := gin.Default()

	r.LoadHTMLGlob("web/templates/*.html")

	// r.Static("/static", "./web/static")
	// Frontend routes
	r.GET("/register", func(c *gin.Context) {
		c.HTML(http.StatusOK, "register.html", gin.H{})
	})
	r.GET("/login", func(c *gin.Context) {
		c.HTML(http.StatusOK, "login.html", gin.H{})
	})
	r.GET("/", func(c *gin.Context) { // Simple home page
		c.HTML(http.StatusOK, "base.html", gin.H{"title": "Welcome Home"})
	})

	// use error (handler) middleware
	var errorMiddleware middlewares.ErrorMiddleware
	r.Use(errorMiddleware.GetHandlerFunc())

	// no auth
	publicAPI := r.Group("/api")
	{
		auth := publicAPI.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
		}
	}

	// requrires auth (bearer token)
	protectedAPI := r.Group("/api")

	var authMiddleware middlewares.AuthMiddleware
	protectedAPI.Use(authMiddleware.GetHandlerFunc())

	{
		auth := protectedAPI.Group("/auth")
		{
			auth.GET("/self", authHandler.Self)
		}

		users := protectedAPI.Group("/users")
		{
			users.GET("/", userHandler.GetAllUsers)
			users.POST("/", userHandler.CreateUser)
			users.GET("/:id", userHandler.GetUserByID)
			users.POST(":id/balance", userHandler.IncrementBalance)
			users.PUT("/:id", userHandler.UpdateUser)
		}

		protectedAPI.GET("/profile", func(c *gin.Context) {
			username, _ := c.Get("username")
			email, _ := c.Get("email")
			c.JSON(http.StatusOK, gin.H{"username": username, "email": email, "message": "Authenticated user profile"})
		})
	}

	return r
}
