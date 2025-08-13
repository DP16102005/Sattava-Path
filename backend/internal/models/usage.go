package models

import (
	"time"

	"gorm.io/gorm"
)

type UsageLog struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	UserID    uint           `json:"user_id" gorm:"not null"`
	AppName   string         `json:"app_name" gorm:"not null"`
	Date      time.Time      `json:"date" gorm:"type:date;not null"`
	Duration  int            `json:"duration"` // in minutes
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	// Relationships
	User User `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

type UsageRequest struct {
	AppName  string `json:"app_name" binding:"required"`
	Duration int    `json:"duration" binding:"required,min=0"`
}

type DailyUsage struct {
	Date     time.Time `json:"date"`
	AppName  string    `json:"app_name"`
	Duration int       `json:"duration"`
}

type UsageStats struct {
	TotalScreenTime int            `json:"total_screen_time"`
	AverageDaily    int            `json:"average_daily"`
	MostUsedApp     string         `json:"most_used_app"`
	AppBreakdown    map[string]int `json:"app_breakdown"`
	WeeklyTrend     []DailyUsage   `json:"weekly_trend"`
}
