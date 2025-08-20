package services

import (
	"errors"
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

func (s *UserService) GetAllUsersPaginated(page, limit int64, query string) (any, pagination.Pagination, error) {
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

	return filteredUser, pagination, nil
}
