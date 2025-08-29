package models

import (
	"time"

	"grocademy/internal/pkg/string_array"

	"gorm.io/gorm"
)

type Course struct {
	ID             uint                     `gorm:"primaryKey" json:"id" faker:"-"`
	CreatedAt      time.Time                `json:"created_at" faker:"-"`
	UpdatedAt      time.Time                `json:"updated_at" faker:"-"`
	DeletedAt      gorm.DeletedAt           `gorm:"index" json:"deleted_at,omitempty" swaggerignore:"true" faker:"-"`
	Title          string                   `json:"title" gorm:"unique;not null" faker:"sentence"`
	Description    string                   `json:"description" gorm:"type:text" faker:"paragraph"`
	Instructor     string                   `json:"instructor" faker:"name"`
	Topics         string_array.StringArray `json:"topics" gorm:"type:text[]" faker:"topics"`
	Price          float64                  `json:"price" faker:"amount"`
	ThumbnailImage string                   `json:"thumbnail_image" faker:"thumbnail"`
}
