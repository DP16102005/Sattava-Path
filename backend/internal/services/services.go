package services

import (
	"digital-wellbeing-backend/internal/config"

	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

type Services struct {
	Auth    *AuthService
	User    *UserService
	Goal    *GoalService
	CheckIn *CheckInService
	Chat    *ChatService
	Usage   *UsageService
}

func NewServices(db *gorm.DB, redis *redis.Client, cfg *config.Config) *Services {
	return &Services{
		Auth:    NewAuthService(db, cfg),
		User:    NewUserService(db),
		Goal:    NewGoalService(db),
		CheckIn: NewCheckInService(db),
		Chat:    NewChatService(db, cfg),
		Usage:   NewUsageService(db),
	}
}
