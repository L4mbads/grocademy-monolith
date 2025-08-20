package handlers

import (
	"net/http"
	"time"

	"grocademy/internal/services" // Assuming services package contains AuthServicer

	"github.com/gin-gonic/gin"
)

type RegisterRequest struct {
	Username  string `json:"username" binding:"required"`
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required,min=6"`
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
// @Router /register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err := h.AuthService.RegisterUser(req.Username, req.Email, req.Password, req.FirstName, req.LastName)
	if err != nil {
		if err.Error() == "username already taken" || err.Error() == "email already registered" {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register user: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully"})
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
// @Router /login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	username, token, err := h.AuthService.LoginUser(req.Identifier, req.Password)
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
