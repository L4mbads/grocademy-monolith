package services

import (
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"grocademy/internal/db/models"
	"grocademy/internal/pkg/pagination"

	"gorm.io/gorm"
)

// ModuleServicer defines the interface for module-related operations.
type ModuleServicer interface {
	CreateModule(courseID uint, title, description string, order int, pdf *multipart.FileHeader, video *multipart.FileHeader) (*models.Module, error)
	GetModuleByID(id uint) (*models.Module, error)
	GetAllModulesByCourseID(courseID uint, page, limit int64, query string) (any, pagination.Pagination, error)
	UpdateModule(id uint, updates map[string]interface{}, pdf *multipart.FileHeader, video *multipart.FileHeader) (*models.Module, error)
	DeleteModule(id uint) error
	ReorderModules(courseID uint, moduleOrders []models.Module) error // Expects a slice of Module with ID and Order
}

// ModuleService implements ModuleServicer.
type ModuleService struct {
	DB *gorm.DB
}

// NewModuleService creates a new ModuleService.
func NewModuleService(db *gorm.DB) *ModuleService {
	return &ModuleService{DB: db}
}

// CreateModule creates a new module for a given course, handling file uploads.
func (s *ModuleService) CreateModule(
	courseID uint, title, description string, order int,
	pdf *multipart.FileHeader, video *multipart.FileHeader,
) (*models.Module, error) {
	// Check if the course exists
	var course models.Course
	if err := s.DB.First(&course, courseID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("course not found")
		}
		return nil, fmt.Errorf("database error checking course: %w", err)
	}

	var pdfPath string
	if pdf != nil {
		path, err := saveContentFile(pdf, "pdfs")
		if err != nil {
			return nil, err
		}
		pdfPath = path
	}

	var videoPath string
	if video != nil {
		path, err := saveContentFile(video, "videos")
		if err != nil {
			return nil, err
		}
		videoPath = path
	}

	module := models.Module{
		CourseID:    courseID,
		Title:       title,
		Description: description,
		Order:       order,
		PDFPath:     pdfPath,
		VideoPath:   videoPath,
	}

	if result := s.DB.Create(&module); result.Error != nil {
		return nil, fmt.Errorf("failed to create module in DB: %w", result.Error)
	}

	return &module, nil
}

// GetModuleByID retrieves a module by its ID.
func (s *ModuleService) GetModuleByID(id uint) (*models.Module, error) {
	var module models.Module
	result := s.DB.First(&module, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("module not found")
		}
		return nil, fmt.Errorf("database error finding module: %w", result.Error)
	}
	return &module, nil
}

// GetAllModulesByCourseID retrieves all modules for a specific course with pagination and search.
func (s *ModuleService) GetAllModulesByCourseID(courseID uint, page, limit int64, query string) (any, pagination.Pagination, error) {
	var modules []models.Module
	searchableColumns := []string{"title", "description"} // Columns to search within modules

	// Filter by CourseID
	dbQuery := s.DB.Where("course_id = ?", courseID).Order("\"order\" ASC") // Order by "order" column

	filteredModules, pagination, err := pagination.Paginate(
		dbQuery.Model(&models.Module{}),
		&modules,
		page,
		limit,
		searchableColumns,
		query,
	)
	if err != nil {
		return nil, pagination, err
	}
	return filteredModules, pagination, nil
}

// UpdateModule updates an existing module, handling partial updates and optional file updates.
func (s *ModuleService) UpdateModule(id uint, updates map[string]interface{}, pdf *multipart.FileHeader, video *multipart.FileHeader) (*models.Module, error) {
	var module models.Module
	result := s.DB.First(&module, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("module not found")
		}
		return nil, fmt.Errorf("database error finding module: %w", result.Error)
	}

	// Handle PDF file update
	if pdf != nil {
		newPDFPath, err := saveContentFile(pdf, "pdfs")
		if err != nil {
			return nil, err
		}
		updates["PDFPath"] = newPDFPath
		// Optionally, delete old PDF file
		if module.PDFPath != "" {
			if err := os.Remove(module.PDFPath); err != nil {
				fmt.Printf("Warning: Failed to delete old PDF file %s: %v\n", module.PDFPath, err)
			}
		}
	} else if _, ok := updates["pdf_content"]; ok && updates["pdf_content"] == nil { // Check if client explicitly sent null to clear
		updates["PDFPath"] = ""
		if module.PDFPath != "" {
			if err := os.Remove(module.PDFPath); err != nil {
				fmt.Printf("Warning: Failed to delete old PDF file %s when clearing: %v\n", module.PDFPath, err)
			}
		}
	}

	// Handle Video file update
	if video != nil {
		newVideoPath, err := saveContentFile(video, "videos")
		if err != nil {
			return nil, err
		}
		updates["VideoPath"] = newVideoPath
		// Optionally, delete old video file
		if module.VideoPath != "" {
			if err := os.Remove(module.VideoPath); err != nil {
				fmt.Printf("Warning: Failed to delete old video file %s: %v\n", module.VideoPath, err)
			}
		}
	} else if _, ok := updates["video_content"]; ok && updates["video_content"] == nil { // Check if client explicitly sent null to clear
		updates["VideoPath"] = ""
		if module.VideoPath != "" {
			if err := os.Remove(module.VideoPath); err != nil {
				fmt.Printf("Warning: Failed to delete old video file %s when clearing: %v\n", module.VideoPath, err)
			}
		}
	}

	// Apply other updates
	if err := s.DB.Model(&module).Updates(updates).Error; err != nil {
		return nil, fmt.Errorf("failed to update module: %w", err)
	}

	return &module, nil
}

// DeleteModule performs a soft delete on a module record.
func (s *ModuleService) DeleteModule(id uint) error {
	var module models.Module
	result := s.DB.First(&module, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return errors.New("module not found")
		}
		return fmt.Errorf("database error finding module: %w", result.Error)
	}

	if deleteResult := s.DB.Delete(&module); deleteResult.Error != nil {
		return fmt.Errorf("failed to delete module: %w", deleteResult.Error)
	}

	// Optionally, delete associated files on soft delete
	if module.PDFPath != "" {
		if err := os.Remove(module.PDFPath); err != nil {
			fmt.Printf("Warning: Failed to delete PDF file %s on soft delete: %v\n", module.PDFPath, err)
		}
	}
	if module.VideoPath != "" {
		if err := os.Remove(module.VideoPath); err != nil {
			fmt.Printf("Warning: Failed to delete video file %s on soft delete: %v\n", module.VideoPath, err)
		}
	}

	return nil
}

// ReorderModules updates the order of multiple modules within a specific course.
func (s *ModuleService) ReorderModules(courseID uint, moduleOrders []models.Module) error {
	tx := s.DB.Begin() // Start a transaction for atomicity
	if tx.Error != nil {
		return fmt.Errorf("failed to begin transaction: %w", tx.Error)
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 1. Validate that all module IDs belong to the specified course
	var existingModules []models.Module
	moduleIDs := make([]uint, len(moduleOrders))
	orderMap := make(map[uint]int) // Map for quick lookup of new orders
	for i, mo := range moduleOrders {
		moduleIDs[i] = mo.ID
		orderMap[mo.ID] = mo.Order
	}

	// Fetch existing modules to verify ownership
	if err := tx.Where("id IN (?) AND course_id = ?", moduleIDs, courseID).Find(&existingModules).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("database error verifying module ownership: %w", err)
	}
	if len(existingModules) != len(moduleOrders) {
		tx.Rollback()
		return errors.New("some module IDs do not belong to the specified course or are invalid")
	}

	// 2. Perform updates in a loop
	for _, moduleData := range moduleOrders {
		if err := tx.Model(&models.Module{}).Where("id = ?", moduleData.ID).Update("order", moduleData.Order).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to update order for module %d: %w", moduleData.ID, err)
		}
	}

	return tx.Commit().Error // Commit the transaction
}

// saveContentFile is a helper function to store uploaded PDF/Video files.
func saveContentFile(file *multipart.FileHeader, subDir string) (string, error) {
	ext := filepath.Ext(file.Filename)
	filename := fmt.Sprintf("%d-%s%s", time.Now().UnixNano(), strconv.FormatInt(time.Now().Unix(), 10), ext) // More unique name
	savePath := filepath.Join("storage", subDir, filename)

	storageDir := filepath.Dir(savePath)
	if _, err := os.Stat(storageDir); os.IsNotExist(err) {
		if err := os.MkdirAll(storageDir, 0755); err != nil {
			return "", fmt.Errorf("failed to create storage directory %s: %w", storageDir, err)
		}
	}

	src, err := file.Open()
	if err != nil {
		return "", fmt.Errorf("failed to open uploaded file: %w", err)
	}
	defer src.Close()

	dst, err := os.Create(savePath)
	if err != nil {
		return "", fmt.Errorf("failed to create destination file %s: %w", savePath, err)
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		return "", fmt.Errorf("failed to save file: %w", err)
	}
	return savePath, nil
}
