package services

import (
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"

	"grocademy/internal/db/models"
	"grocademy/internal/dto"
	"grocademy/internal/pkg/pagination"
	"grocademy/internal/storage"

	"gorm.io/gorm"
)

type CourseServicer interface {
	CreateCourse(title, description, instructor string, topics []string, price float64, thumbnail *multipart.FileHeader) (*models.Course, error)
	GetCourseByID(userID, courseID uint) (*models.Course, int64, bool, error)
	GetMyCourses(userID uint, page, limit int64) (*[]dto.MyCourseResponse, pagination.Pagination, error)
	GetAllCoursesPaginated(page, limit int64, query string) (*[]map[string]interface{}, pagination.Pagination, error)
	UpdateCourse(id uint, updates map[string]interface{}, thumbnail *multipart.FileHeader) (*models.Course, error)
	DeleteCourse(id uint) error
	BuyCourse(userID uint, courseID uint) (float64, uint, error)
}

type CourseService struct {
	DB    *gorm.DB
	Cloud storage.CloudStorage
}

func NewCourseService(db *gorm.DB, cloud storage.CloudStorage) *CourseService {
	return &CourseService{DB: db, Cloud: cloud}
}

func (s *CourseService) CreateCourse(
	title, description, instructor string,
	topics []string,
	price float64,
	thumbnail *multipart.FileHeader,
) (*models.Course, error) {
	var thumbnailPath string

	if thumbnail != nil {
		URL, err := s.saveThumbnail(thumbnail, title)
		if err != nil {
			return nil, err
		}

		thumbnailPath = URL
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

func (s *CourseService) GetCourseByID(userID, courseID uint) (*models.Course, int64, bool, error) {
	var course models.Course
	result := s.DB.First(&course, courseID)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, 0, false, errors.New("course not found")
		}
		return nil, 0, false, fmt.Errorf("database error finding course: %w", result.Error)
	}
	var totalModules int64
	s.DB.Model(&models.Module{}).Where("course_id = ?", courseID).Count(&totalModules)

	var purchased bool
	err := s.DB.Model(&models.Enrollment{}).Select("count(*) > 0").Where("user_id = ? AND course_id = ?", userID, courseID).Find(&purchased).Error
	if err != nil {
		return nil, 0, false, err
	}

	return &course, totalModules, purchased, nil
}

func (s *CourseService) GetAllCoursesPaginated(page, limit int64, query string) (*[]map[string]interface{}, pagination.Pagination, error) {
	// var courses []models.Course
	// searchableColumns := []string{"title", "instructor", "topics"}

	// filteredCourses, pagination, err := pagination.Paginate(
	// 	s.DB.Model(&models.Course{}),
	// 	&courses,
	// 	page,
	// 	limit,
	// 	searchableColumns,
	// 	query,
	// )

	// assertedCourses, _ := filteredCourses.(*[]models.Course)

	// if err != nil {
	// 	return nil, pagination, err
	// }
	// return assertedCourses, pagination, nil
	// A custom struct for the query result, including the module count.
	var results []struct {
		models.Course
		TotalModules int64
	}

	// Build the base query for both filtering and counting.
	dbQuery := s.DB.Model(&models.Course{}).
		Select("courses.*, count(modules.id) as total_modules").
		Joins("left join modules on modules.course_id = courses.id").
		Group("courses.id")

	// Apply search filtering
	if query != "" {
		searchQuery := ""
		searchableColumns := []string{"title", "description", "instructor", "topics"}
		args := make([]interface{}, len(searchableColumns))
		for i, col := range searchableColumns {
			if col == "topics" {
				searchQuery += "EXISTS (SELECT 1 FROM unnest(courses.topics) AS t WHERE t ILIKE ?)"
			} else {
				searchQuery += fmt.Sprintf("courses.%s ILIKE ?", col)
			}
			if i < len(searchableColumns)-1 {
				searchQuery += " OR "
			}
			args[i] = fmt.Sprintf("%%%s%%", query)
		}
		dbQuery = dbQuery.Where(searchQuery, args...)
	}

	// Count total items for pagination first
	var totalItems int64
	dbQuery.Count(&totalItems)

	// Paginate the query and execute
	_, pagination, err := pagination.Paginate(dbQuery, &results, page, limit, nil, "")
	if err != nil {
		return nil, pagination, err
	}

	// Transfer data to a final response format
	var coursesWithCount []map[string]interface{}
	for _, res := range results {
		courseMap := map[string]interface{}{
			"id":              res.ID,
			"title":           res.Title,
			"description":     res.Description,
			"instructor":      res.Instructor,
			"topics":          res.Topics,
			"price":           res.Price,
			"thumbnail_image": res.ThumbnailImage,
			"created_at":      res.CreatedAt,
			"updated_at":      res.UpdatedAt,
			"deleted_at":      res.DeletedAt,
			"total_modules":   res.TotalModules, // ADDED
		}
		coursesWithCount = append(coursesWithCount, courseMap)
	}

	return &coursesWithCount, pagination, nil
}

func (s *CourseService) GetMyCourses(userID uint, page, limit int64) (*[]dto.MyCourseResponse, pagination.Pagination, error) {
	// A custom struct to hold the data from our complex join query.
	var results []struct {
		models.Course
		models.Enrollment
	}

	// First, count total enrolled courses for the pagination metadata.
	var totalItems int64
	s.DB.Model(&models.Enrollment{}).Where("user_id = ?", userID).Count(&totalItems)

	// Build the main query for enrolled courses, joining with `enrollments` table.
	s.DB.Model(&models.Course{}).
		Select("courses.*, enrollments.*").
		Joins("INNER JOIN enrollments ON enrollments.course_id = courses.id").
		Where("enrollments.user_id = ?", userID).
		Find(&results)

	// _, pagination, err := pagination.Paginate(query, &results, page, limit, nil, "")
	// if err != nil {
	// 	return nil, pagination, fmt.Errorf("failed to paginate user's courses: %w", err)
	// }

	// var myCourses []map[string]interface{}
	myCourses := []dto.MyCourseResponse{}
	for _, enrolledCourse := range results {
		// Get total modules for the course
		var totalModules int64
		s.DB.Model(&models.Module{}).Where("course_id = ?", enrolledCourse.ID).Count(&totalModules)

		// Get total completed modules for this user and course
		var completedModules int64
		s.DB.Model(&models.ModuleProgress{}).
			Joins("JOIN modules ON modules.id = module_progress.module_id").
			Where("module_progress.user_id = ? AND modules.course_id = ? AND module_progress.is_completed = ?", userID, enrolledCourse.ID, true).
			Count(&completedModules)

		// Calculate progress percentage
		progressPercentage := 0.0
		if totalModules > 0 {
			progressPercentage = float64(completedModules) / float64(totalModules) * 100
		}

		// courseData := map[string]interface{}{
		// 	"id":                  enrolledCourse.ID,
		// 	"title":               enrolledCourse.Title,
		// 	"description":         enrolledCourse.Description,
		// 	"instructor":          enrolledCourse.Instructor,
		// 	"topics":              enrolledCourse.Topics,
		// 	"price":               enrolledCourse.Price,
		// 	"thumbnail_image":     enrolledCourse.ThumbnailImage,
		// 	"total_modules":       totalModules,
		// 	"progress_percentage": progressPercentage,
		// 	"purchased_at":        enrolledCourse.PurchaseDate,
		// }
		myCourses = append(myCourses, dto.MyCourseResponse{
			Enrollment:         enrolledCourse.Enrollment,
			Course:             enrolledCourse.Course,
			ProgressPercentage: progressPercentage,
		})
	}
	pagination := pagination.Pagination{}

	pagination.CurrentPage = page
	pagination.TotalItems = totalItems
	pagination.TotalPages = int64(totalItems / limit)
	if totalItems%limit != 0 {
		pagination.TotalPages++
	}

	return &myCourses, pagination, nil
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
		newPath, err := s.saveThumbnail(thumbnail, course.Title)
		if err != nil {
			return nil, err
		}
		updates["ThumbnailImage"] = newPath

		// Optionally, delete the old thumbnail file
		if course.ThumbnailImage != "" {
			if err := os.Remove(course.ThumbnailImage); err != nil {
				fmt.Printf("Warning: Failed to delete old thumbnail image %s: %v\n", course.ThumbnailImage, err)
			}
		}
	} else if _, ok := updates["thumbnail_image"]; ok && updates["thumbnail_image"] == nil {
		// If thumbnail_image was explicitly sent as null/empty string, clear the path
		updates["ThumbnailImage"] = ""
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
func (s *CourseService) saveThumbnail(thumbnail *multipart.FileHeader, title string) (string, error) {
	// ext := filepath.Ext(thumbnail.Filename)
	// filename := fmt.Sprintf("%d-%s%s", time.Now().UnixNano(), strings.ReplaceAll(strings.ToLower(title), " ", "-"), ext)
	savePath := filepath.Join("course", "thumbnail", thumbnail.Filename)

	// create directory if exisn't.
	storageDir := filepath.Dir(savePath)
	if _, err := os.Stat(storageDir); os.IsNotExist(err) {
		if err := os.MkdirAll(storageDir, 0755); err != nil {
			return "", fmt.Errorf("failed to create storage directory: %w", err)
		}
	}

	// save to local

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

	URL, err := s.Cloud.UploadFile(thumbnail, savePath)
	if err != nil {
		return "", fmt.Errorf("failed to upload to cloud: %w", err)
	}
	return URL, nil
	// ext := filepath.Ext(thumbnail.Filename)
	// filename := fmt.Sprintf("%d-%s%s", time.Now().UnixNano(), strings.ReplaceAll(strings.ToLower(title), " ", "-"), ext)
	// savePath := filepath.Join("storage", "images", filename)

	// storageDir := filepath.Dir(savePath)
	// if _, err := os.Stat(storageDir); os.IsNotExist(err) {
	// 	if err := os.MkdirAll(storageDir, 0755); err != nil {
	// 		return "", fmt.Errorf("failed to create storage directory: %w", err)
	// 	}
	// }

	// src, err := thumbnail.Open()
	// if err != nil {
	// 	return "", fmt.Errorf("failed to open uploaded file: %w", err)
	// }
	// defer src.Close()

	// dst, err := os.Create(savePath)
	// if err != nil {
	// 	return "", fmt.Errorf("failed to create destination file: %w", err)
	// }
	// defer dst.Close()

	// if _, err := io.Copy(dst, src); err != nil {
	// 	return "", fmt.Errorf("failed to save file: %w", err)
	// }

	// URL, err := s.Cloud.UploadFile(thumbnail, savePath)
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to upload to cloud: %w", err)
	// }
	// return savePath, nil
}

func (s *CourseService) BuyCourse(userID uint, courseID uint) (float64, uint, error) {
	tx := s.DB.Begin() // Start a transaction for atomicity
	if tx.Error != nil {
		return 0, 0, fmt.Errorf("failed to begin transaction: %w", tx.Error)
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 1. Check if the course exists and get its price.
	var course models.Course
	if err := tx.First(&course, courseID).Error; err != nil {
		tx.Rollback()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, 0, errors.New("course not found")
		}
		return 0, 0, fmt.Errorf("database error checking course: %w", err)
	}

	// 2. Check if the user already purchased the course.
	var existingEnrollment models.Enrollment
	if err := tx.Where("user_id = ? AND course_id = ?", userID, courseID).First(&existingEnrollment).Error; err == nil {
		tx.Rollback()
		return 0, 0, errors.New("user has already purchased this course")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		tx.Rollback()
		return 0, 0, fmt.Errorf("database error checking existing enrollment: %w", err)
	}

	// 3. Check user balance.
	var user models.User
	if err := tx.First(&user, userID).Error; err != nil {
		tx.Rollback()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, 0, errors.New("user not found")
		}
		return 0, 0, fmt.Errorf("database error checking user balance: %w", err)
	}

	if user.Balance < course.Price {
		tx.Rollback()
		return user.Balance, 0, errors.New("insufficient balance")
	}

	// 4. Reduce user balance.
	user.Balance -= course.Price
	if err := tx.Save(&user).Error; err != nil {
		tx.Rollback()
		return user.Balance, 0, fmt.Errorf("failed to reduce user balance: %w", err)
	}

	// 5. Create a new enrollment entry.
	enrollment := models.Enrollment{
		UserID:   userID,
		CourseID: courseID,
	}
	if err := tx.Create(&enrollment).Error; err != nil {
		tx.Rollback()
		return user.Balance, 0, fmt.Errorf("failed to create enrollment: %w", err)
	}

	return user.Balance, enrollment.TransactionID, tx.Commit().Error // Commit the transaction
}
