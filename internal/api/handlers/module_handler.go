package handlers

import (
	"errors"
	"fmt"
	"mime/multipart"
	"net/http"
	"strconv"

	"grocademy/internal/db/models"
	"grocademy/internal/services"

	"github.com/gin-gonic/gin"
)

// CreateModuleRequest defines the form data for creating a module.
type CreateModuleRequest struct {
	Title        string                `form:"title" binding:"required"`
	Description  string                `form:"description" binding:"required"`
	PDFContent   *multipart.FileHeader `form:"pdf_content"`
	VideoContent *multipart.FileHeader `form:"video_content"`
}

// UpdateModuleRequest defines the form data for updating a module.
type UpdateModuleRequest struct {
	Title        string                `form:"title,omitempty"`
	Description  string                `form:"description,omitempty"`
	PDFContent   *multipart.FileHeader `form:"pdf_content,omitempty"`
	VideoContent *multipart.FileHeader `form:"video_content,omitempty"`
	// Consider adding fields to explicitly clear PDF/Video content if needed
	ClearPDF   bool `form:"clear_pdf,omitempty"`   // Example for clearing content
	ClearVideo bool `form:"clear_video,omitempty"` // Example for clearing content
}

// ReorderModulesRequest defines the request body for reordering modules.
type ReorderModulesRequest struct {
	ModuleOrder []struct {
		ID    uint `json:"id" binding:"required"`
		Order int  `json:"order" binding:"required"`
	} `json:"module_order" binding:"required"`
}

type CompleteModuleRequest struct {
	IsCompleted *bool `json:"is_completed" binding:"required"`
}

// ModuleHandler handles module-related API requests.
type ModuleHandler struct {
	ModuleService services.ModuleServicer
}

// NewModuleHandler creates a new ModuleHandler.
func NewModuleHandler(moduleService services.ModuleServicer) *ModuleHandler {
	return &ModuleHandler{ModuleService: moduleService}
}

// CreateModule godoc
// @Summary Create a new module for a specific course
// @Description Create a new module from multipart form data for the given course ID
// @Tags modules
// @Accept  multipart/form-data
// @Produce  json
// @Param courseId path int true "Course ID"
// @Param title formData string true "Module title"
// @Param description formData string true "Module description"
// @Param order formData int true "Module order within the course"
// @Param pdf_content formData file false "PDF file for module content"
// @Param video_content formData file false "Video file for module content"
// @Success 201 {object} models.Module
// @Failure 400 {object} map[string]string "Invalid input"
// @Failure 404 {object} map[string]string "Course not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Security Bearer
// @Router /courses/{courseId}/modules [post]
func (h *ModuleHandler) CreateModule(c *gin.Context) {
	courseIDStr := c.Param("id")
	courseID, err := strconv.ParseUint(courseIDStr, 10, 64)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, errors.New("invalid course ID"))
		return
	}

	var req CreateModuleRequest
	if err := c.ShouldBind(&req); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	newModule, err := h.ModuleService.CreateModule(
		uint(courseID),
		req.Title,
		req.Description,
		req.PDFContent,
		req.VideoContent,
	)
	if err != nil {
		if err.Error() == "course not found" {
			c.AbortWithError(http.StatusNotFound, err)
			return
		}
		c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("failed to create module: %v", err))
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"message": "module created",
		"data":    newModule,
	})
}

// GetAllModulesByCourseID godoc
// @Summary Get all modules for a specific course with pagination and search
// @Description Retrieve a list of all modules for a given course, with optional pagination and search parameters
// @Tags modules
// @Produce  json
// @Param courseId path int true "Course ID"
// @Param page query int false "Page number (default 1)"
// @Param limit query int false "Items per page (default 15)"
// @Param q query string false "Search query"
// @Success 200 {object} []models.Module
// @Failure 400 {object} map[string]string "Invalid input"
// @Failure 500 {object} map[string]string "Internal server error"
// @Security Bearer
// @Router /courses/{courseId}/modules [get]
func (h *ModuleHandler) GetAllModulesByCourseID(c *gin.Context) {
	courseIDStr := c.Param("id")
	courseID, err := strconv.ParseUint(courseIDStr, 10, 64)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, errors.New("invalid course ID"))
		return
	}

	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "15")

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
	limit = min(limit, 50)

	userID, _ := c.Get("id")

	paginatedModules, progressMap, pagination, err := h.ModuleService.GetAllModulesByCourseID(uint(courseID), userID.(uint), int64(page), int64(limit))
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("failed to retrieve modules: %v", err))
		return
	}

	var enrichedModules []map[string]interface{}
	for _, module := range *paginatedModules {
		enrichedModule := map[string]interface{}{
			"id":            module.ID,
			"course_id":     module.CourseID,
			"title":         module.Title,
			"description":   module.Description,
			"order":         module.Order,
			"pdf_content":   module.PDFPath,
			"video_content": module.VideoPath,
			"created_at":    module.CreatedAt,
			"updated_at":    module.UpdatedAt,
			"is_completed":  (*progressMap)[module.ID],
		}
		enrichedModules = append(enrichedModules, enrichedModule)
	}

	c.JSON(http.StatusOK, gin.H{
		"status":     "success",
		"message":    "modules queried",
		"data":       enrichedModules,
		"pagination": pagination,
	})
}

// GetModuleByID godoc
// @Summary Get a module by ID
// @Description Retrieve a single module by its ID
// @Tags modules
// @Produce  json
// @Param id path int true "Module ID"
// @Success 200 {object} models.Module
// @Failure 400 {object} map[string]string "Invalid module ID"
// @Failure 404 {object} map[string]string "Module not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Security Bearer
// @Router /modules/{id} [get]
func (h *ModuleHandler) GetModuleByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, errors.New("invalid course ID"))
		return
	}

	userID, _ := c.Get("id")

	module, completion, err := h.ModuleService.GetModuleByID(uint(id), userID.(uint))
	if err != nil {
		if err.Error() == "module not found" {
			c.AbortWithError(http.StatusNotFound, err)
			return
		}
		c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("failed to retrieve module: %v", err))
		return
	}

	enrichedModule := map[string]interface{}{
		"id":            module.ID,
		"course_id":     module.CourseID,
		"title":         module.Title,
		"description":   module.Description,
		"order":         module.Order,
		"pdf_content":   module.PDFPath,
		"video_content": module.VideoPath,
		"created_at":    module.CreatedAt,
		"updated_at":    module.UpdatedAt,
		"is_completed":  completion,
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "module found",
		"data":    enrichedModule,
	})
}

// UpdateModule godoc
// @Summary Update a module's data
// @Description Update specified fields of a module by ID, with optional file uploads
// @Tags modules
// @Accept  multipart/form-data
// @Produce  json
// @Param id path int true "Module ID"
// @Param title formData string false "Module title"
// @Param description formData string false "Module description"
// @Param order formData int false "Module order within the course"
// @Param pdf_content formData file false "New PDF file for module content"
// @Param video_content formData file false "New Video file for module content"
// @Param clear_pdf formData boolean false "Set to true to clear existing PDF content"
// @Param clear_video formData boolean false "Set to true to clear existing Video content"
// @Success 200 {object} models.Module "Updated module object"
// @Failure 400 {object} map[string]string "Invalid input or no fields to update"
// @Failure 404 {object} map[string]string "Module not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Security Bearer
// @Router /modules/{id} [put]
func (h *ModuleHandler) UpdateModule(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, errors.New("invalid course ID"))
		return
	}

	var req UpdateModuleRequest
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

	// Handle explicit clearing of PDF/Video content
	if req.ClearPDF {
		updates["pdf_content"] = nil
	}
	if req.ClearVideo {
		updates["video_content"] = nil
	}

	if len(updates) == 0 && req.PDFContent == nil && req.VideoContent == nil && !req.ClearPDF && !req.ClearVideo {
		c.AbortWithError(http.StatusBadRequest, errors.New("no fields to update provided"))
		return
	}

	updatedModule, err := h.ModuleService.UpdateModule(uint(id), updates, req.PDFContent, req.VideoContent)
	if err != nil {
		if err.Error() == "module not found" {
			c.AbortWithError(http.StatusNotFound, err)
			return
		}
		c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("failed to update module: %v", err))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "module updated",
		"data":    updatedModule,
	})
}

// DeleteModule godoc
// @Summary Delete a module
// @Description Deletes a module record by ID (soft delete)
// @Tags modules
// @Produce  json
// @Param id path int true "Module ID"
// @Success 204 "Module deleted successfully"
// @Failure 400 {object} map[string]string "Invalid module ID"
// @Failure 404 {object} map[string]string "Module not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Security Bearer
// @Router /modules/{id} [delete]
func (h *ModuleHandler) DeleteModule(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, errors.New("invalid course ID"))
		return
	}

	if err := h.ModuleService.DeleteModule(uint(id)); err != nil {
		if err.Error() == "module not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Invalid Module ID"})
			return
		}
		c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("failed to delete module: %v", err))
		return
	}

	c.Status(http.StatusNoContent)
}

// ReorderModules godoc
// @Summary Reorder modules within a course
// @Description Update the order of multiple modules for a specific course
// @Tags modules
// @Accept  json
// @Produce  json
// @Param courseId path int true "Course ID"
// @Param module_order body ReorderModulesRequest true "List of module IDs and their new orders"
// @Success 200 {object} map[string]string "message: Modules reordered successfully"
// @Failure 400 {object} map[string]string "Invalid input or modules not belonging to course"
// @Failure 500 {object} map[string]string "Internal server error"
// @Security Bearer
// @Router /courses/{courseId}/modules/reorder [patch]
func (h *ModuleHandler) ReorderModules(c *gin.Context) {
	courseIDStr := c.Param("id")
	courseID, err := strconv.ParseUint(courseIDStr, 10, 64)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, errors.New("invalid course ID"))
		return
	}

	var req ReorderModulesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	// Convert request data to models.Module slice for the service
	modulesToReorder := make([]models.Module, len(req.ModuleOrder))
	for i, item := range req.ModuleOrder {
		modulesToReorder[i] = models.Module{
			ID:    item.ID,
			Order: item.Order,
		}
	}

	if err := h.ModuleService.ReorderModules(uint(courseID), modulesToReorder); err != nil {
		c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("failed to reorder modules: %v", err))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Modules reordered successfully",
		"data":    req.ModuleOrder,
	})
}

// CompleteModuleByID godoc
// @Summary Get a module by ID
// @Description Retrieve a single module by its ID
// @Tags modules
// @Produce  json
// @Param id path int true "Module ID"
// @Success 200 {object} models.Module
// @Failure 400 {object} map[string]string "Invalid module ID"
// @Failure 404 {object} map[string]string "Module not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Security Bearer
// @Router /modules/{id}/complete [patch]
func (h *ModuleHandler) CompleteModuleByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, errors.New("invalid course ID"))
		return
	}
	var req CompleteModuleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	userID, _ := c.Get("id")

	total, completed, percentage, latestCompletion, err := h.ModuleService.CompleteModuleByID(uint(id), userID.(uint), *req.IsCompleted)
	if err != nil {
		if err.Error() == "module not found" {
			c.AbortWithError(http.StatusNotFound, err)
			return
		}
		c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("failed to retrieve module: %v", err))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "module found",
		"data": gin.H{
			"module_id":    id,
			"is_completed": *req.IsCompleted,
			"course_progress": gin.H{
				"total_modules":     total,
				"completed_modules": completed,
				"percentage":        percentage,
				"latest_completion": latestCompletion,
			},
		},
	})
}
