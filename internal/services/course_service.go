package services

import (
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"grocademy/internal/db/models"
	"grocademy/internal/pkg/pagination"

	"gorm.io/gorm"
)

type CourseServicer interface {
	CreateCourse(title, description, instructor string, topics []string, price float64, thumbnail *multipart.FileHeader) (*models.Course, error)
	GetCourseByID(id uint) (*models.Course, error)
	GetAllCoursesPaginated(page, limit int64, query string) (*[]models.Course, pagination.Pagination, error)
	UpdateCourse(id uint, updates map[string]interface{}, thumbnail *multipart.FileHeader) (*models.Course, error)
	DeleteCourse(id uint) error
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
		Title:          title,
		Description:    description,
		Instructor:     instructor,
		Topics:         topics,
		Price:          price,
		ThumbnailImage: thumbnailPath,
	}

	if result := s.DB.Create(&course); result.Error != nil {
		return nil, fmt.Errorf("failed to create course in DB: %w", result.Error)
	}

	return &course, nil
}

func (s *CourseService) GetCourseByID(id uint) (*models.Course, error) {
	var course models.Course
	result := s.DB.First(&course, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("course not found")
		}
		return nil, fmt.Errorf("database error finding course: %w", result.Error)
	}
	return &course, nil
}

func (s *CourseService) GetAllCoursesPaginated(page, limit int64, query string) (*[]models.Course, pagination.Pagination, error) {
	var courses []models.Course
	searchableColumns := []string{"title", "instructor", "topics"}

	filteredCourses, pagination, err := pagination.Paginate(
		s.DB.Model(&models.Course{}),
		&courses,
		page,
		limit,
		searchableColumns,
		query,
	)

	assertedCourses, _ := filteredCourses.(*[]models.Course)

	if err != nil {
		return nil, pagination, err
	}
	return assertedCourses, pagination, nil
}

func (s *CourseService) UpdateCourse(id uint, updates map[string]interface{}, thumbnail *multipart.FileHeader) (*models.Course, error) {
	var course models.Course
	result := s.DB.First(&course, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("course not found")
		}
		return nil, fmt.Errorf("database error finding course: %w", result.Error)
	}

	// Handle thumbnail image update if provided
	if thumbnail != nil {
		// Save new thumbnail and update path
		newPath, err := saveThumbnail(thumbnail, course.Title)
		if err != nil {
			return nil, err
		}
		updates["ThumbnailPath"] = newPath

		// Optionally, delete the old thumbnail file
		if course.ThumbnailImage != "" {
			if err := os.Remove(course.ThumbnailImage); err != nil {
				fmt.Printf("Warning: Failed to delete old thumbnail image %s: %v\n", course.ThumbnailImage, err)
			}
		}
	} else if _, ok := updates["thumbnail_image"]; ok && updates["thumbnail_image"] == nil {
		// If thumbnail_image was explicitly sent as null/empty string, clear the path
		updates["ThumbnailPath"] = ""
		if course.ThumbnailImage != "" {
			if err := os.Remove(course.ThumbnailImage); err != nil {
				fmt.Printf("Warning: Failed to delete old thumbnail image %s when clearing: %v\n", course.ThumbnailImage, err)
			}
		}
	}

	fmt.Printf("ini tipe nya %T\n", updates["Topics"])

	if err := s.DB.Model(&course).Updates(updates).Error; err != nil {
		return nil, fmt.Errorf("failed to update course: %w", err)
	}

	return &course, nil
}

func (s *CourseService) DeleteCourse(id uint) error {
	var course models.Course
	result := s.DB.First(&course, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return errors.New("course not found")
		}
		return fmt.Errorf("database error finding course: %w", result.Error)
	}

	if deleteResult := s.DB.Delete(&course); deleteResult.Error != nil {
		return fmt.Errorf("failed to delete course: %w", deleteResult.Error)
	}

	// Optionally, delete the thumbnail file from storage on soft delete
	if course.ThumbnailImage != "" {
		if err := os.Remove(course.ThumbnailImage); err != nil {
			fmt.Printf("Warning: Failed to delete thumbnail image %s on soft delete: %v\n", course.ThumbnailImage, err)
		}
	}

	return nil
}

// saveThumbnail is a helper function to store the uploaded image.
func saveThumbnail(thumbnail *multipart.FileHeader, title string) (string, error) {
	ext := filepath.Ext(thumbnail.Filename)
	filename := fmt.Sprintf("%d-%s%s", time.Now().UnixNano(), strings.ReplaceAll(strings.ToLower(title), " ", "-"), ext)
	savePath := filepath.Join("storage", "images", filename)

	storageDir := filepath.Dir(savePath)
	if _, err := os.Stat(storageDir); os.IsNotExist(err) {
		if err := os.MkdirAll(storageDir, 0755); err != nil {
			return "", fmt.Errorf("failed to create storage directory: %w", err)
		}
	}

	src, err := thumbnail.Open()
	if err != nil {
		return "", fmt.Errorf("failed to open uploaded file: %w", err)
	}
	defer src.Close()

	dst, err := os.Create(savePath)
	if err != nil {
		return "", fmt.Errorf("failed to create destination file: %w", err)
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		return "", fmt.Errorf("failed to save file: %w", err)
	}
	return savePath, nil
}
