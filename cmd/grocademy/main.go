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

	db.Init()
	gormDB := db.GetDB()

	userService := services.NewUserService(gormDB)

	userHandler := handlers.NewUserHandler(userService)

	router := api.SetupRouter(userHandler)

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on :%s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
