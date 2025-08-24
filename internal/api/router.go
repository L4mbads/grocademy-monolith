package api

import (
	"log"
	"net/http"
	"os"
	"time"

	"grocademy/internal/api/handlers"
	"grocademy/internal/api/middlewares"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	docs "grocademy/api"

	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type StartableRouter interface {
	Start()
}

type GinRouterWrapper struct {
	ginEngine *gin.Engine
}

func NewRouter(
	userHandler *handlers.UserHandler,
	authHandler *handlers.AuthHandler,
	courseHandler *handlers.CourseHandler,
	moduleHandler *handlers.ModuleHandler,
) GinRouterWrapper {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	docs.SwaggerInfo.BasePath = "/api"

	// CORS config
	r.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	r.LoadHTMLFiles(
		"web/templates/components/navbar.html",
		"web/templates/components/search_bar.html",
		"web/templates/components/paging.html",
		"web/templates/layout.html",
		"web/templates/dashboard.html",
		"web/templates/browse_courses.html",
		"web/templates/course_modules.html",
		"web/templates/course.html",
		"web/templates/register.html",
		"web/templates/login.html")
	r.Static("/static", "./web/static")

	// use error (handler) middleware
	var errorMiddleware middlewares.ErrorMiddleware
	r.Use(errorMiddleware.GetHandlerFunc())

	// Public FE routes
	r.GET("/register", func(c *gin.Context) {
		c.HTML(http.StatusOK, "register.html", gin.H{})
	})
	r.GET("/login", func(c *gin.Context) {
		c.HTML(http.StatusOK, "login.html", gin.H{})
	})

	// Protected FE routes
	authWebMiddleware := middlewares.NewAuthWebMiddleware()
	authenticatedWeb := r.Group("")
	authenticatedWeb.Use(authWebMiddleware.GetHandlerFunc())
	{
		authenticatedWeb.GET("", func(c *gin.Context) {
			c.Redirect(http.StatusMovedPermanently, "/dashboard")
		})
		authenticatedWeb.GET("/dashboard", func(c *gin.Context) {
			c.HTML(http.StatusOK, "dashboard.html", gin.H{})
		})
		authenticatedWeb.GET("/courses", func(c *gin.Context) {
			c.HTML(http.StatusOK, "browse_courses.html", gin.H{"title": "Browse Course", "route": "api/courses"})
		})
		authenticatedWeb.GET("/my-courses", func(c *gin.Context) {
			c.HTML(http.StatusOK, "browse_courses.html", gin.H{"title": "My Course", "route": "api/courses/my-courses"})
		})
		authenticatedWeb.GET("/courses/:id", func(c *gin.Context) {
			c.HTML(http.StatusOK, "course.html", gin.H{})
		})
		authenticatedWeb.GET("/courses/:id/modules", func(c *gin.Context) {
			c.HTML(http.StatusOK, "course_modules.html", gin.H{})
		})
	}

	r.GET("/health", func(c *gin.Context) { c.JSON(200, gin.H{"status": "ok"}) })

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

	authAPIMiddleware := middlewares.NewAuthAPIMiddleware()
	protectedAPI.Use(authAPIMiddleware.GetHandlerFunc())

	adminMiddleware := middlewares.NewAdminMiddleware()
	{
		auth := protectedAPI.Group("/auth")
		{
			auth.GET("/self", authHandler.Self)
		}

		users := protectedAPI.Group("/users")
		users.Use(adminMiddleware.GetHandlerFunc())
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
			courses.GET("", courseHandler.GetAllCourses)
			courses.GET("/my-courses", courseHandler.GetMyCourses)
			courses.POST("/:id/buy", courseHandler.BuyCourse)
			courses.GET("/:id", courseHandler.GetCourseByID)

			protectedCourses := courses.Group("")
			protectedCourses.Use(adminMiddleware.GetHandlerFunc())

			protectedCourses.POST("", courseHandler.CreateCourse)
			protectedCourses.PUT("/:id", courseHandler.UpdateCourse)
			protectedCourses.DELETE("/:id", courseHandler.DeleteCourse)

			modulesByCourse := courses.Group("/:id/modules")
			{
				modulesByCourse.GET("", moduleHandler.GetAllModulesByCourseID)

				protectedModulesByCourse := modulesByCourse.Group("")
				protectedModulesByCourse.Use(adminMiddleware.GetHandlerFunc())

				protectedModulesByCourse.POST("", moduleHandler.CreateModule)
				protectedModulesByCourse.PATCH("/reorder", moduleHandler.ReorderModules)
			}
		}

		modules := protectedAPI.Group("/modules")
		{
			modules.GET("/:id", moduleHandler.GetModuleByID)
			modules.PATCH("/:id/complete", moduleHandler.CompleteModuleByID)

			protectedModules := modules.Group("")
			protectedModules.Use(adminMiddleware.GetHandlerFunc())
			protectedModules.PUT("/:id", moduleHandler.UpdateModule)
			protectedModules.DELETE("/:id", moduleHandler.DeleteModule)
		}
	}

	r.GET("/docs/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	r.RemoveExtraSlash = true

	return GinRouterWrapper{ginEngine: r}
}

func (g GinRouterWrapper) Start() {
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on :%s", port)
	if err := g.ginEngine.Run(":" + port); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
