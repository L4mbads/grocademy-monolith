package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username  string `json:"username" gorm:"unique;not null"`
	Email     string `json:"email" gorm:"unique;not null"`
	Password  string `json:"password" gorm:"not null"`
	FirstName string `json:"first_name" gorm:"not null"`
	LastName  string `json:"last_name" gorm:"not null"`
	Balance   int    `json:"balance" gorm:"not null"`
}
