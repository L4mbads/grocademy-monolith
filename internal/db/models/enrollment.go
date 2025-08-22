package models

import (
	"time"

	"gorm.io/gorm"
)

type Enrollment struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	UserID       uint           `json:"user_id" gorm:"not null;uniqueIndex:uq_user_course"`
	User         User           `json:"-"` // GORM association, json:"-" hides it from JSON output
	CourseID     uint           `json:"course_id" gorm:"not null;uniqueIndex:uq_user_course"`
	Course       Course         `json:"-"` // GORM association
	PurchaseDate time.Time      `json:"purchase_date" gorm:"default:CURRENT_TIMESTAMP"`
}
