package handlers

import (
	"fmt"
	"mime/multipart"
	"net/http"

	"grocademy/internal/services"

	"github.com/gin-gonic/gin"
)

type CreateCourseRequest struct {
	Title          string                `form:"title" binding:"required"`
	Description    string                `form:"description" binding:"required"`
	Instructor     string                `form:"instructor" binding:"required"`
	Topics         []string              `form:"topics" binding:"required"` // Bind as a single string, then split
	Price          float64               `form:"price" binding:"required"`
	ThumbnailImage *multipart.FileHeader `form:"thumbnail_image"` // The binary image file
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
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to create course: %v", err)})
		return
	}

	c.JSON(http.StatusCreated, newCourse)
}
