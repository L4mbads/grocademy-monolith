package db

import (
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"grocademy/internal/auth"
	"grocademy/internal/db/models"
)

var DB *gorm.DB

func Init() {

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=Asia/Jakarta",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_SSLMODE"),
	)

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second,   // Slow SQL threshold
			LogLevel:                  logger.Silent, // Log level
			IgnoreRecordNotFoundError: true,          // Ignore ErrRecordNotFound error for logger
			ParameterizedQueries:      true,          // Log parameterized queries, changes from `sql` to `sql?`
			Colorful:                  false,         // Disable color
		},
	)

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	err = DB.AutoMigrate(
		&models.User{},
		&models.Course{},
		&models.Module{},
		&models.Enrollment{},
		&models.ModuleProgress{},
	)
	if err != nil {
		log.Fatalf("Failed to auto migrate database: %v", err)
	}

	log.Println("Database connection established and migrations applied (if any).")

	createDefaultAdmin(DB)
}

func GetDB() *gorm.DB {
	return DB
}

func createDefaultAdmin(db *gorm.DB) {
	adminUsername := "admin"
	adminEmail := "admin@example.com"
	adminPassword := "admin123"
	adminFirstName := "admin"
	adminLastName := "admin"
	adminBalance := 9999999999.0

	var adminUser models.User
	result := db.Where("email = ?", adminEmail).First(&adminUser)

	if result.Error != nil && result.Error == gorm.ErrRecordNotFound {
		hashedPassword, err := auth.HashPassword(adminPassword)
		if err != nil {
			log.Fatalf("Failed to hash admin password: %v", err)
		}

		newAdmin := models.User{
			Username:  adminUsername,
			Email:     adminEmail,
			Password:  hashedPassword,
			FirstName: adminFirstName,
			LastName:  adminLastName,
			Balance:   adminBalance,
		}

		if createResult := db.Create(&newAdmin); createResult.Error != nil {
			log.Fatalf("Failed to create default admin user: %v", createResult.Error)
		}
		log.Println("Default admin user created successfully.")
	} else if result.Error != nil {
		log.Fatalf("Database error while checking for admin user: %v", result.Error)
	} else {
		log.Println("Default admin user already exists.")
	}
}
