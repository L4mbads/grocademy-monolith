package models

import (
	"time"

	"gorm.io/gorm"
)

type Course struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	Title         string         `json:"title" gorm:"not null"`
	Description   string         `json:"description" gorm:"type:text"`
	Instructor    string         `json:"instructor"`
	Topics        string         `json:"topics" gorm:"type:text"`
	Price         float64        `json:"price"`
	ThumbnailPath string         `json:"thumbnail_path"`
}
