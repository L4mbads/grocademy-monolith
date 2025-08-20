package services

import (
	"errors"
	"fmt"
	"grocademy/internal/auth"
	"grocademy/internal/db/models"

	"gorm.io/gorm"
)

type AuthServicer interface {
	RegisterUser(username, email, password, firstName, lastName string) (*models.User, error)
	LoginUser(email, password string) (string, string, error)
	GetCurrentUser(username string) (*models.User, error)
}

type AuthService struct {
	DB *gorm.DB
}

func NewAuthService(db *gorm.DB) *AuthService {
	return &AuthService{DB: db}
}

func (s *AuthService) RegisterUser(username, email, password, firstName, lastName string) (*models.User, error) {

	var existingUser models.User

	if err := s.DB.Where("username = ?", username).Or("email = ?", email).First(&existingUser).Error; err == nil {
		if existingUser.Username == username {
			return nil, errors.New("username already taken")
		}
		if existingUser.Email == email {
			return nil, errors.New("email already registered")
		}
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("database error checking existing user: %w", err)
	}

	hashedPassword, err := auth.HashPassword(password)
	if err != nil {
		return nil, errors.New("failed to hash password")
	}

	user := &models.User{
		Username:  username,
		Email:     email,
		Password:  hashedPassword,
		FirstName: firstName,
		LastName:  lastName,
		Balance:   0,
	}

	if result := s.DB.Create(user); result.Error != nil {
		return nil, fmt.Errorf("failed to register user: %w", result.Error)
	}

	return user, nil
}

func (s *AuthService) LoginUser(identifier, password string) (string, string, error) {

	var user models.User

	result := s.DB.Where("email = ?", identifier).Or("username = ?", identifier).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return "", "", errors.New("invalid credentials")
		}

		return "", "", fmt.Errorf("database error during login: %w", result.Error)
	}

	if !auth.CheckPasswordHash(password, user.Password) {
		return "", "", errors.New("invalid credentials")
	}

	token, err := auth.GenerateJWT(user.Username, user.Email)
	if err != nil {
		return "", "", errors.New("failed to generate token")
	}

	return user.Username, token, nil
}

func (s *AuthService) GetCurrentUser(username string) (*models.User, error) {
	var user models.User

	result := s.DB.Where("username = ?", username).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("invalid credentials")
		}

		return nil, fmt.Errorf("database error: %w", result.Error)
	}

	return &user, nil

}
