package api

import (
	"net/http"
	"time"

	"grocademy/internal/api/handlers"
	"grocademy/internal/api/middlewares"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// SetupRouter configures the Gin router with API routes.
func SetupRouter(
	userHandler *handlers.UserHandler,
	authHandler *handlers.AuthHandler,
	courseHandler *handlers.CourseHandler,
) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	// CORS config
	r.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	r.GET("/health", func(c *gin.Context) { c.JSON(200, gin.H{"status": "ok"}) })
	r.LoadHTMLGlob("web/templates/*.html")

	// r.Static("/static", "./web/static")
	// Frontend routes
	r.GET("/register", func(c *gin.Context) {
		c.HTML(http.StatusOK, "register.html", gin.H{})
	})
	r.GET("/login", func(c *gin.Context) {
		c.HTML(http.StatusOK, "login.html", gin.H{})
	})
	r.GET("", func(c *gin.Context) { // Simple home page
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
			users.GET("", userHandler.GetAllUsers)
			users.POST("", userHandler.CreateUser)

			user := users.Group("/:id")
			{
				user.GET("", userHandler.GetUserByID)
				user.PUT("", userHandler.UpdateUser)
				user.DELETE("", userHandler.DeleteUser)
				user.POST("/balance", userHandler.IncrementBalance)
			}
		}

		courses := protectedAPI.Group("/courses")
		{
			courses.POST("", courseHandler.CreateCourse)
			courses.GET("", courseHandler.GetAllCourses)
			courses.GET("/:id", courseHandler.GetCourseByID)
			courses.PUT("/:id", courseHandler.UpdateCourse)
			courses.DELETE("/:id", courseHandler.DeleteCourse)
		}
	}

	r.RemoveExtraSlash = true

	return r
}
