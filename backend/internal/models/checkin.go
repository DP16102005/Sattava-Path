package models

import (
	"time"

	"gorm.io/gorm"
)

type CheckIn struct {
	ID         uint           `json:"id" gorm:"primaryKey"`
	UserID     uint           `json:"user_id" gorm:"not null"`
	Date       time.Time      `json:"date" gorm:"type:date;not null"`
	Mood       int            `json:"mood" gorm:"not null;check:mood >= 1 AND mood <= 5"`
	ScreenTime int            `json:"screen_time"` // in minutes
	Notes      string         `json:"notes"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `json:"-" gorm:"index"`

	// Relationships
	User User `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

type CheckInRequest struct {
	Mood       int    `json:"mood" binding:"required,min=1,max=5"`
	ScreenTime int    `json:"screen_time" binding:"required,min=0"`
	Notes      string `json:"notes"`
}

type MoodStats struct {
	AverageMood      float64     `json:"average_mood"`
	MoodDistribution map[int]int `json:"mood_distribution"`
	Streak           int         `json:"streak"`
	TotalCheckIns    int         `json:"total_check_ins"`
}
