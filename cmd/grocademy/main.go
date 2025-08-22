package main

import (
	"grocademy/internal/api"
	"grocademy/internal/api/handlers"
	"grocademy/internal/db"
	"grocademy/internal/services"
)

// @title           Grocademy API
// @version         1.0
// @description     Fachriza Ahmad Setiyono - Seleksi Labpro 3
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /api

// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/

func main() {

	// Initialize database
	db.Init()
	gormDB := db.GetDB()

	// Initialize services
	userService := services.NewUserService(gormDB)
	authService := services.NewAuthService(gormDB)
	courseService := services.NewCourseService(gormDB)
	moduleService := services.NewModuleService(gormDB)

	// Initialize handlers
	userHandler := handlers.NewUserHandler(userService)
	authHandler := handlers.NewAuthHandler(authService)
	courseHandler := handlers.NewCourseHandler(courseService)
	moduleHandler := handlers.NewModuleHandler(moduleService)

	router := api.NewRouter(
		userHandler,
		authHandler,
		courseHandler,
		moduleHandler,
	)
	router.Start()

}
