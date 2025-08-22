package models

import (
	"time"

	"gorm.io/gorm"
)

type ModuleProgress struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty" swaggerignore:"true"`
	UserID      uint           `json:"user_id" gorm:"not null;uniqueIndex:uq_user_module"`
	User        User           `json:"-"` // GORM association
	ModuleID    uint           `json:"module_id" gorm:"not null;uniqueIndex:uq_user_module"`
	Module      Module         `json:"-"` // GORM association
	IsCompleted bool           `json:"is_completed" gorm:"default:false"`
	CompletedAt *time.Time     `json:"completed_at"` // Pointer to time.Time to allow NULL in DB
}
