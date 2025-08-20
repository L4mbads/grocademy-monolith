// cmd/web/main.go
package main

import (
	"log"
	"os"

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
	userService := services.NewUserService(gormDB) // Existing user service
	authService := services.NewAuthService(gormDB) // NEW: Auth service

	// Initialize handlers
	userHandler := handlers.NewUserHandler(userService)
	authHandler := handlers.NewAuthHandler(authService) // NEW: Auth handler

	// Setup router - now takes both handlers
	router := api.SetupRouter(userHandler, authHandler) // UPDATED

	// Get port from environment variables, default to 8080
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on :%s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
