package services

import (
	"errors"
	"fmt"
	"grocademy/internal/auth"
	"grocademy/internal/db/models"
	"grocademy/internal/pkg/pagination"

	"gorm.io/gorm"
)

type UserService struct {
	DB *gorm.DB
}

func NewUserService(db *gorm.DB) *UserService {
	return &UserService{DB: db}
}

func (s *UserService) CreateUser(user *models.User) error {
	result := s.DB.Create(user)
	return result.Error
}

func (s *UserService) GetUserByID(id uint) (*models.User, error) {
	var user models.User
	result := s.DB.First(&user, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, result.Error
	}
	return &user, nil
}

func (s *UserService) GetUsers() ([]models.User, error) {
	var users []models.User
	result := s.DB.Find(&users)
	if result.Error != nil {
		return nil, result.Error
	}
	return users, nil
}

func (s *UserService) GetAllUsersPaginated(page, limit int64, query string) (*[]models.User, pagination.Pagination, error) {
	var users []models.User
	searchableColumns := []string{"username", "email", "first_name", "last_name"}

	filteredUser, pagination, err := pagination.Paginate(
		s.DB.Model(&models.User{}),
		&users,
		page,
		limit,
		searchableColumns,
		query,
	)
	if err != nil {
		return nil, pagination, err
	}

	assertedUser := filteredUser.(*[]models.User)

	return assertedUser, pagination, nil
}

func (s *UserService) UpdateUser(id uint, updates map[string]interface{}) (*models.User, error) {
	var user models.User
	result := s.DB.First(&user, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, fmt.Errorf("database error finding user: %w", result.Error)
	}

	if user.Username == "admin" {
		return nil, errors.New("update on admin user is prohibited")
	}

	if updates["Password"] != "" {
		hashedPassword, err := auth.HashPassword(fmt.Sprint(updates["Password"]))
		if err != nil {
			return nil, errors.New("failed to hash password")
		}
		updates["Password"] = hashedPassword
	}

	if err := s.DB.Model(&user).Updates(updates).Error; err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	return &user, nil
}

func (s *UserService) IncrementUserBalance(id uint, increment float64) (*models.User, error) {
	var user models.User
	result := s.DB.First(&user, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, fmt.Errorf("database error finding user: %w", result.Error)
	}

	if err := s.DB.Model(&user).Update("balance", user.Balance+increment).Error; err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	return &user, nil
}

func (s *UserService) DeleteUser(id uint) error {
	var user models.User
	result := s.DB.First(&user, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return errors.New("user not found")
		}
		return fmt.Errorf("database error finding user: %w", result.Error)
	}
	if user.Username == "admin" {
		return errors.New("deletion on admin user is prohibited")
	}

	// GORM's soft delete
	if deleteResult := s.DB.Delete(&user); deleteResult.Error != nil {
		return fmt.Errorf("failed to delete user: %w", deleteResult.Error)
	}

	return nil
}
