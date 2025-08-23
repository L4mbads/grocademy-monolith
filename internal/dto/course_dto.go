package dto

import (
	"grocademy/internal/db/models"
)

type MyCourseResponse struct {
	models.Course
	models.Enrollment
	ProgressPercentage float64 `json:"progress_percentage"`
}
