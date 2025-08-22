package main

import (
	"grocademy/internal/api"
	"grocademy/internal/api/handlers"
	"grocademy/internal/db"
	"grocademy/internal/services"
)

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
