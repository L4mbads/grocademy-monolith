package handlers

import (
	"errors"
	"fmt"
	"mime/multipart"
	"net/http"
	"strconv"

	"grocademy/internal/db/models"
	"grocademy/internal/pkg/string_array"
	"grocademy/internal/services"

	"github.com/gin-gonic/gin"
)

type CourseResponse struct {
	models.Course
	TotalModules int64 `json:"total_modules"` // NEW
}

type CreateCourseRequest struct {
	Title          string                `form:"title" binding:"required"`
	Description    string                `form:"description" binding:"required"`
	Instructor     string                `form:"instructor" binding:"required"`
	Topics         []string              `form:"topics" binding:"required"` // Bind as a single string, then split
	Price          float64               `form:"price" binding:"required"`
	ThumbnailImage *multipart.FileHeader `form:"thumbnail_image"` // The binary image file
}

// For partial updates, fields are optional.
type UpdateCourseRequest struct {
	Title          string                `form:"title,omitempty"`
	Description    string                `form:"description,omitempty"`
	Instructor     string                `form:"instructor,omitempty"`
	Topics         []string              `form:"topics,omitempty"`
	Price          *float64              `form:"price,omitempty"`           // Use pointer for float64 to distinguish 0 from unset
	ThumbnailImage *multipart.FileHeader `form:"thumbnail_image,omitempty"` // Optional file upload
}

type CourseHandler struct {
	CourseService services.CourseServicer
}

func NewCourseHandler(courseService services.CourseServicer) *CourseHandler {
	return &CourseHandler{CourseService: courseService}
}

// CreateCourse handles the POST request to create a new course.
// @Summary Create a new course
// @Description Create a new course from multipart form data
// @Tags courses
// @Accept  multipart/form-data
// @Produce  json
// @Param title formData string true "Course title"
// @Param description formData string true "Course description"
// @Param instructor formData string true "Course instructor"
// @Param topics formData string true "List of topics"
// @Param price formData number true "Course price"
// @Param thumbnail_image formData file false "Thumbnail image file"
// @Success 201 {object} models.Course
// @Failure 400 {object} map[string]string "Invalid input"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /courses [post]
func (h *CourseHandler) CreateCourse(c *gin.Context) {
	var req CreateCourseRequest
	if err := c.ShouldBind(&req); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	newCourse, err := h.CourseService.CreateCourse(
		req.Title,
		req.Description,
		req.Instructor,
		req.Topics,
		req.Price,
		req.ThumbnailImage,
	)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("failed to create course: %v", err))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Query success",
		"data":    newCourse,
	})
}

// GetCourseByID godoc (NEW HANDLER)
// @Summary Get a course by ID
// @Description Retrieve a single course by its ID
// @Tags courses
// @Produce  json
// @Param id path int true "Course ID"
// @Success 200 {object} models.Course
// @Failure 400 {object} map[string]string "Invalid course ID"
// @Failure 404 {object} map[string]string "Course not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /courses/{id} [get]
func (h *CourseHandler) GetCourseByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, errors.New("invalid course ID"))
		return
	}

	course, totalModules, err := h.CourseService.GetCourseByID(uint(id))
	if err != nil {
		if err.Error() == "course not found" {
			c.AbortWithError(http.StatusNotFound, err)
			return
		}
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	response := CourseResponse{
		Course:       *course,
		TotalModules: totalModules,
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Query success",
		"data":    response,
	})
}

// GetAllCourses godoc (NEW HANDLER)
// @Summary Get all courses with pagination and search
// @Description Retrieve a list of all courses with optional pagination and search parameters
// @Tags courses
// @Produce  json
// @Param page query int false "Page number (default 1)"
// @Param limit query int false "Items per page (default 10)"
// @Param q query string false "Search query"
// @Success 200 {object} pagination.PaginatedResponse{data=[]models.Course}
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /courses [get]
func (h *CourseHandler) GetAllCourses(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "10")
	query := c.DefaultQuery("q", "")

	page, err := strconv.ParseInt(pageStr, 10, 64)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, errors.New("invalid page number"))
		return
	}
	limit, err := strconv.ParseInt(limitStr, 10, 64)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, errors.New("invalid limit number"))
		return
	}

	paginatedCourses, pagination, err := h.CourseService.GetAllCoursesPaginated(int64(page), int64(limit), query)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":     "success",
		"message":    "Query success",
		"data":       paginatedCourses,
		"pagination": pagination,
	})
}

// UpdateCourse godoc (NEW HANDLER)
// @Summary Update a course's data
// @Description Update specified fields of a course by ID, with optional thumbnail upload
// @Tags courses
// @Accept  multipart/form-data
// @Produce  json
// @Param id path int true "Course ID"
// @Param title formData string false "Course title"
// @Param description formData string false "Course description"
// @Param instructor formData string false "Course instructor"
// @Param topics formData string false "Comma-separated list of topics"
// @Param price formData number false "Course price"
// @Param thumbnail_image formData file false "New thumbnail image file"
// @Success 200 {object} models.Course "Updated course object"
// @Failure 400 {object} map[string]string "Invalid input or no fields to update"
// @Failure 404 {object} map[string]string "Course not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /courses/{id} [put]
func (h *CourseHandler) UpdateCourse(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, errors.New("invalid course ID"))
		return
	}

	var req UpdateCourseRequest
	// Use c.ShouldBind to handle multipart/form-data
	if err := c.ShouldBind(&req); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	updates := make(map[string]interface{})
	if req.Title != "" {
		updates["Title"] = req.Title
	}
	if req.Description != "" {
		updates["Description"] = req.Description
	}
	if req.Instructor != "" {
		updates["Instructor"] = req.Instructor
	}
	if req.Price != nil {
		updates["Price"] = *req.Price
	}

	if req.Topics != nil {
		updates["Topics"] = string_array.StringArray(req.Topics)
	}

	if len(updates) == 0 && req.ThumbnailImage == nil {
		c.AbortWithError(http.StatusBadRequest, errors.New("no fields to update provided"))
		return
	}

	updatedCourse, err := h.CourseService.UpdateCourse(uint(id), updates, req.ThumbnailImage)
	if err != nil {
		if err.Error() == "course not found" {
			c.AbortWithError(http.StatusNotFound, err)
			return
		}
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "user updated",
		"data":    updatedCourse,
	})
}

// DeleteCourse godoc (NEW HANDLER)
// @Summary Delete a course
// @Description Deletes a course record by ID (soft delete)
// @Tags courses
// @Produce  json
// @Param id path int true "Course ID"
// @Success 204 "Course deleted successfully"
// @Failure 400 {object} map[string]string "Invalid course ID"
// @Failure 404 {object} map[string]string "Course not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /courses/{id} [delete]
func (h *CourseHandler) DeleteCourse(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.AbortWithError(http.StatusNotFound, errors.New("invalid course ID"))
		return
	}

	if err := h.CourseService.DeleteCourse(uint(id)); err != nil {
		if err.Error() == "course not found" {
			c.AbortWithError(http.StatusNotFound, err)
			return
		}
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.Status(http.StatusNoContent)
}
