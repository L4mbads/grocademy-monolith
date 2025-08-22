package models

import (
	"time"

	"gorm.io/gorm"
)

type Module struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty" swaggerignore:"true"`
	CourseID    uint           `json:"course_id" gorm:"not null"` // Foreign key to Course
	Course      Course         `json:"-"`                         // GORM association
	Title       string         `json:"title" gorm:"not null"`
	Description string         `json:"description" gorm:"type:text"`
	Order       int            `json:"order" gorm:"not null"` // Module order within the course
	PDFPath     string         `json:"pdf_content"`           // Path to the stored PDF file
	VideoPath   string         `json:"video_content"`         // Path to the stored video file
}
