package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"grocademy/internal/db/models"
	"grocademy/internal/services"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	UserService *services.UserService
}

// NewUserHandler creates a new UserHandler
func NewUserHandler(userService *services.UserService) *UserHandler {
	return &UserHandler{UserService: userService}
}

// CreateUser godoc
// @Summary Create a new user
// @Description Create a new user with the input payload
// @Tags users
// @Accept  json
// @Produce  json
// @Param user body models.User true "User object to be created"
// @Success 200 {object} models.User
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /users [post]
func (h *UserHandler) CreateUser(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	if err := h.UserService.CreateUser(&user); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"message": "User created",
		"data":    user,
	})
}

// GetUserByID godoc
// @Summary Get a user by ID
// @Description Get a single user by their ID
// @Tags users
// @Produce  json
// @Param id path int true "User ID"
// @Success 200 {object} models.User
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /users/{id} [get]
func (h *UserHandler) GetUserByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, errors.New("invalid user ID"))
		return
	}

	user, err := h.UserService.GetUserByID(uint(id))
	if err != nil {
		if err.Error() == "user not found" {
			c.AbortWithError(http.StatusNotFound, err)
			return
		}
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "User found",
		"data":    user,
	})
}

// GetAllUsers godoc
// @Summary Get all users with pagination and search
// @Description Retrieve a list of all users with optional pagination and search parameters
// @Tags users
// @Produce  json
// @Param page query int false "Page number (default 1)"
// @Param limit query int false "Items per page (default 10)"
// @Param q query string false "Search query"
// @Success 200 {object} pagination.PaginatedResponse{data=[]models.User}
// @Failure 500 {object} map[string]string
// @Router /users [get]
func (h *UserHandler) GetAllUsers(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "15")
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

	page = max(page, 0)
	limit = max(50, min(limit, 0))

	paginatedUsers, pagination, err := h.UserService.GetAllUsersPaginated(int64(page), int64(limit), query)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":     "success",
		"message":    "Query success",
		"data":       paginatedUsers,
		"pagination": pagination,
	})
}
