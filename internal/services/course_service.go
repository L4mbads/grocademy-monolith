package services

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"grocademy/internal/db/models"

	"gorm.io/gorm"
)

// CourseServicer defines the interface for course-related operations.
type CourseServicer interface {
	CreateCourse(title, description, instructor string, topics []string, price float64, thumbnail *multipart.FileHeader) (*models.Course, error)
}

// CourseService implements CourseServicer.
type CourseService struct {
	DB *gorm.DB
}

// NewCourseService creates a new CourseService.
func NewCourseService(db *gorm.DB) *CourseService {
	return &CourseService{DB: db}
}

// CreateCourse handles saving the course details and the thumbnail image.
func (s *CourseService) CreateCourse(
	title, description, instructor string,
	topics []string,
	price float64,
	thumbnail *multipart.FileHeader,
) (*models.Course, error) {
	var thumbnailPath string

	if thumbnail != nil {
		// 1. Define a unique file name and path for storage.
		ext := filepath.Ext(thumbnail.Filename)
		filename := fmt.Sprintf("%d-%s%s", time.Now().UnixNano(), strings.ReplaceAll(strings.ToLower(title), " ", "-"), ext)
		savePath := filepath.Join("storage", "images", filename)

		// 2. Create the storage directory if it doesn't exist.
		storageDir := filepath.Dir(savePath)
		if _, err := os.Stat(storageDir); os.IsNotExist(err) {
			if err := os.MkdirAll(storageDir, 0755); err != nil {
				return nil, fmt.Errorf("failed to create storage directory: %w", err)
			}
		}

		// 3. Save the uploaded file to the local storage.
		src, err := thumbnail.Open()
		if err != nil {
			return nil, fmt.Errorf("failed to open uploaded file: %w", err)
		}
		defer src.Close()

		dst, err := os.Create(savePath)
		if err != nil {
			return nil, fmt.Errorf("failed to create destination file: %w", err)
		}
		defer dst.Close()

		if _, err := io.Copy(dst, src); err != nil {
			return nil, fmt.Errorf("failed to save file: %w", err)
		}
		thumbnailPath = savePath // Store the local path in the database
	}

	// 4. Create the Course model and save to the database.
	course := models.Course{
		Title:         title,
		Description:   description,
		Instructor:    instructor,
		Topics:        strings.Join(topics, ","),
		Price:         price,
		ThumbnailPath: thumbnailPath,
	}

	if result := s.DB.Create(&course); result.Error != nil {
		return nil, fmt.Errorf("failed to create course in DB: %w", result.Error)
	}

	return &course, nil
}
