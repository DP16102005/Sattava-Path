package models

import (
	"time"

	"gorm.io/gorm"
)

type Goal struct {
	ID         uint           `json:"id" gorm:"primaryKey"`
	UserID     uint           `json:"user_id" gorm:"not null"`
	AppName    string         `json:"app_name" gorm:"not null"`
	AppIcon    string         `json:"app_icon"`
	DailyLimit int            `json:"daily_limit"` // in minutes
	IsEnabled  bool           `json:"is_enabled" gorm:"default:true"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `json:"-" gorm:"index"`

	// Relationships
	User User `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

type GoalRequest struct {
	AppName    string `json:"app_name" binding:"required"`
	AppIcon    string `json:"app_icon"`
	DailyLimit int    `json:"daily_limit" binding:"required,min=1"`
	IsEnabled  bool   `json:"is_enabled"`
}

type GoalProgress struct {
	Goal            Goal    `json:"goal"`
	UsedToday       int     `json:"used_today"`     // minutes used today
	RemainingTime   int     `json:"remaining_time"` // minutes remaining
	IsOverLimit     bool    `json:"is_over_limit"`
	ProgressPercent float64 `json:"progress_percent"`
}
