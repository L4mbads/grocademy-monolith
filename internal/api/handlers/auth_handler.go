package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	_ "grocademy/internal/db/models"
	"grocademy/internal/services" // Assuming services package contains AuthServicer

	"github.com/gin-gonic/gin"
)

type RegisterRequest struct {
	Username  string `json:"username" binding:"required"`
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required,min=8"`
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
}

type LoginRequest struct {
	Identifier string `json:"identifier" binding:"required"`
	Password   string `json:"password" binding:"required"`
}

type LoginData struct {
	Username string `json:"username"`
	Token    string `json:"token"`
}

type RegisterData struct {
	ID        string `json:"id"`
	Username  string `json:"username"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type AuthHandler struct {
	AuthService services.AuthServicer
}

func NewAuthHandler(authService services.AuthServicer) *AuthHandler {
	return &AuthHandler{AuthService: authService}
}

// Register godoc
// @Summary Register a new user
// @Description Register a new user with username, email, and password
// @Tags auth
// @Accept  json
// @Produce  json
// @Param user body RegisterRequest true "User registration details"
// @Success 201 {object} map[string]string "message: User registered successfully"
// @Failure 400 {object} map[string]string "error: Invalid input"
// @Failure 409 {object} map[string]string "error: Username or email already taken"
// @Failure 500 {object} map[string]string "error: Internal server error"
// @Router /auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	user, err := h.AuthService.RegisterUser(req.Username, req.Email, req.Password, req.FirstName, req.LastName)
	if err != nil {
		if err.Error() == "username already taken" || err.Error() == "email already registered" {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	data := RegisterData{
		ID:        strconv.FormatUint(uint64(user.ID), 10),
		Username:  req.Username,
		FirstName: req.FirstName,
		LastName:  req.LastName,
	}

	c.JSON(
		http.StatusCreated,
		gin.H{
			"status":  "success",
			"message": "User registered successfully",
			"data":    data,
		})
}

// Login godoc
// @Summary Log in a user
// @Description Authenticate a user with email and password, and return a JWT token in an HttpOnly cookie
// @Tags auth
// @Accept  json
// @Produce  json
// @Param credentials body LoginRequest true "User login credentials"
// @Success 200 {object} map[string]string "message: Login successful"
// @Failure 400 {object} map[string]string "error: Invalid input"
// @Failure 401 {object} map[string]string "error: Invalid credentials"
// @Failure 500 {object} map[string]string "error: Internal server error"
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	site := c.DefaultQuery("site", "admin")
	username, token, err := h.AuthService.LoginUser(req.Identifier, req.Password, site)
	if err != nil {
		if err.Error() == "invalid credentials" {
			c.AbortWithError(http.StatusUnauthorized, err)
			return
		}
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	// Set the JWT token as an HttpOnly cookie
	c.SetCookie(
		"jwt_token",              // Cookie name
		token,                    // Cookie value (the JWT token)
		int(time.Hour.Seconds()), // Max-Age: 1 hour (in seconds)
		"/",                      // Path: Available across the entire domain
		"",                       // Domain: Empty means current domain
		false,                    // Secure: Set to true for HTTPS only in production
		true,                     // HttpOnly: Prevent JavaScript access
	)

	data := LoginData{
		Username: username,
		Token:    token,
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Login successful",
		"data":    data,
	})

}

// Self godoc
// @Summary Get current user
// @Description Get currently logged-in user, based on token
// @Tags auth
// @Produce json
// @Success 200 {object} models.User "User"
// @Failure 500 {object} map[string]string "Internal server error"
// @Security Bearer
// @Router /auth/self [get]
func (h *AuthHandler) Self(c *gin.Context) {
	username, _ := c.Get("username")
	user, err := h.AuthService.GetCurrentUser(fmt.Sprint(username))

	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "self",
		"data":    user,
	})
}
