package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        uint           `gorm:"primaryKey" json:"id" faker:"-"`
	CreatedAt time.Time      `json:"created_at" faker:"-"`
	UpdatedAt time.Time      `json:"updated_at" faker:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty" swaggerignore:"true" faker:"-"`
	Username  string         `json:"username" gorm:"unique;not null" faker:"username"`
	Email     string         `json:"email" gorm:"unique;not null"  faker:"email"`
	Password  string         `json:"-" gorm:"not null" faker:"password"`
	FirstName string         `json:"first_name" gorm:"not null"  faker:"first_name"`
	LastName  string         `json:"last_name" gorm:"not null" faker:"last_name"`
	Balance   float64        `json:"balance" gorm:"not null" faker:"amount"`
}
