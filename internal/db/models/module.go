package models

import (
	"time"

	"gorm.io/gorm"
)

type Module struct {
	ID          uint           `gorm:"primaryKey" json:"id" faker:"-"`
	CreatedAt   time.Time      `json:"created_at"  faker:"-"`
	UpdatedAt   time.Time      `json:"updated_at"  faker:"-"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty" swaggerignore:"true" faker:"-"`
	CourseID    uint           `json:"course_id" gorm:"not null" faker:"course_id"` // Foreign key to Course
	Course      Course         `json:"-" faker:"-"`                                 // GORM association
	Title       string         `json:"title" gorm:"not null" faker:"sentence"`
	Description string         `json:"description" gorm:"type:text" faker:"paragraph"`
	Order       int            `json:"order" gorm:"not null" faker:"order"` // Module order within the course
	PDFPath     string         `json:"pdf_content" faker:"pdf_path"`        // Path to the stored PDF file
	VideoPath   string         `json:"video_content" faker:"video_path"`    // Path to the stored video file
}
