package services

import (
	"digital-wellbeing-backend/internal/models"

	"gorm.io/gorm"
)

type UserService struct {
	db *gorm.DB
}

func NewUserService(db *gorm.DB) *UserService {
	return &UserService{db: db}
}

func (s *UserService) GetProfile(userID uint) (*models.UserResponse, error) {
	var user models.User
	if err := s.db.First(&user, userID).Error; err != nil {
		return nil, err
	}

	response := user.ToResponse()
	return &response, nil
}

func (s *UserService) UpdateProfile(userID uint, updates map[string]interface{}) (*models.UserResponse, error) {
	var user models.User
	if err := s.db.First(&user, userID).Error; err != nil {
		return nil, err
	}

	if err := s.db.Model(&user).Updates(updates).Error; err != nil {
		return nil, err
	}

	response := user.ToResponse()
	return &response, nil
}

func (s *UserService) GetAllUsers() ([]models.UserResponse, error) {
	var users []models.User
	if err := s.db.Find(&users).Error; err != nil {
		return nil, err
	}

	var responses []models.UserResponse
	for _, user := range users {
		responses = append(responses, user.ToResponse())
	}

	return responses, nil
}
