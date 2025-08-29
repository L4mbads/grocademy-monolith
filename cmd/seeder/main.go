package main

import (
	"grocademy/internal/db"
	"grocademy/internal/pkg/seeding"
)

func main() {
	// Initialize database
	db.Init()
	gormDB := db.GetDB()

	seeder := seeding.NewSeeder(gormDB)
	seeder.Seed(20, 20, 2)
}
