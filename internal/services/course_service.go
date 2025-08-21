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

type CourseServicer interface {
	CreateCourse(title, description, instructor string, topics []string, price float64, thumbnail *multipart.FileHeader) (*models.Course, error)
}

type CourseService struct {
	DB *gorm.DB
}

func NewCourseService(db *gorm.DB) *CourseService {
	return &CourseService{DB: db}
}

func (s *CourseService) CreateCourse(
	title, description, instructor string,
	topics []string,
	price float64,
	thumbnail *multipart.FileHeader,
) (*models.Course, error) {
	var thumbnailPath string

	if thumbnail != nil {
		ext := filepath.Ext(thumbnail.Filename)
		filename := fmt.Sprintf("%d-%s%s", time.Now().UnixNano(), strings.ReplaceAll(strings.ToLower(title), " ", "-"), ext)
		savePath := filepath.Join("storage", "images", filename)

		// create directory if exisn't.
		storageDir := filepath.Dir(savePath)
		if _, err := os.Stat(storageDir); os.IsNotExist(err) {
			if err := os.MkdirAll(storageDir, 0755); err != nil {
				return nil, fmt.Errorf("failed to create storage directory: %w", err)
			}
		}

		// save to local
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
		thumbnailPath = savePath
	}

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
