package services

import (
	"errors"
	"time"

	"digital-wellbeing-backend/internal/config"
	"digital-wellbeing-backend/internal/models"
	"digital-wellbeing-backend/internal/utils"

	"gorm.io/gorm"
)

type AuthService struct {
	db  *gorm.DB
	cfg *config.Config
}

func NewAuthService(db *gorm.DB, cfg *config.Config) *AuthService {
	return &AuthService{db: db, cfg: cfg}
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type RegisterRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type AuthResponse struct {
	User  models.UserResponse `json:"user"`
	Token string              `json:"token"`
}

func (s *AuthService) Register(req RegisterRequest) (*AuthResponse, error) {
	// Check if user already exists
	var existingUser models.User
	if err := s.db.Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
		return nil, errors.New("user already exists")
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	// Create user
	user := models.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: hashedPassword,
	}

	if err := s.db.Create(&user).Error; err != nil {
		return nil, err
	}

	// Generate JWT
	expiresIn, _ := time.ParseDuration(s.cfg.JWTExpiresIn)
	token, err := utils.GenerateJWT(user.ID, user.Email, user.IsAdmin, s.cfg.JWTSecret, expiresIn)
	if err != nil {
		return nil, err
	}

	return &AuthResponse{
		User:  user.ToResponse(),
		Token: token,
	}, nil
}

func (s *AuthService) Login(req LoginRequest) (*AuthResponse, error) {
	var user models.User
	if err := s.db.Where("email = ?", req.Email).First(&user).Error; err != nil {
		return nil, errors.New("invalid credentials")
	}

	if !utils.CheckPassword(req.Password, user.Password) {
		return nil, errors.New("invalid credentials")
	}

	// Generate JWT
	expiresIn, _ := time.ParseDuration(s.cfg.JWTExpiresIn)
	token, err := utils.GenerateJWT(user.ID, user.Email, user.IsAdmin, s.cfg.JWTSecret, expiresIn)
	if err != nil {
		return nil, err
	}

	return &AuthResponse{
		User:  user.ToResponse(),
		Token: token,
	}, nil
}

func (s *AuthService) GetUserByID(userID uint) (*models.User, error) {
	var user models.User
	if err := s.db.First(&user, userID).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
