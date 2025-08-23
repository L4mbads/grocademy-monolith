package models

import (
	"time"

	"gorm.io/gorm"
)

type Enrollment struct {
	TransactionID uint           `gorm:"primaryKey" json:"transaction_id"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty" swaggerignore:"true"`
	UserID        uint           `json:"user_id" gorm:"not null;uniqueIndex:uq_user_course"`
	User          User           `json:"-"` // GORM association, json:"-" hides it from JSON output
	CourseID      uint           `json:"course_id" gorm:"not null;uniqueIndex:uq_user_course"`
	Course        Course         `json:"-"` // GORM association
	PurchasedAt   time.Time      `json:"purchased_at" gorm:"default:CURRENT_TIMESTAMP"`
}
