package models

import (
	"time"

	"grocademy/internal/pkg/string_array"

	"gorm.io/gorm"
)

type Course struct {
	ID             uint                     `gorm:"primaryKey" json:"id"`
	CreatedAt      time.Time                `json:"created_at"`
	UpdatedAt      time.Time                `json:"updated_at"`
	DeletedAt      gorm.DeletedAt           `gorm:"index" json:"deleted_at,omitempty" swaggerignore:"true"`
	Title          string                   `json:"title" gorm:"unique;not null"`
	Description    string                   `json:"description" gorm:"type:text"`
	Instructor     string                   `json:"instructor"`
	Topics         string_array.StringArray `json:"topics" gorm:"type:text[]"`
	Price          float64                  `json:"price"`
	ThumbnailImage string                   `json:"thumbnail_image"`
}
