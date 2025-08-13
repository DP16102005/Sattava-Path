package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Email     string         `json:"email" gorm:"uniqueIndex;not null"`
	Name      string         `json:"name" gorm:"not null"`
	Password  string         `json:"-" gorm:"not null"`
	Avatar    string         `json:"avatar"`
	IsAdmin   bool           `json:"is_admin" gorm:"default:false"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	// Relationships
	Goals     []Goal     `json:"goals,omitempty"`
	CheckIns  []CheckIn  `json:"check_ins,omitempty"`
	ChatLogs  []ChatLog  `json:"chat_logs,omitempty"`
	UsageLogs []UsageLog `json:"usage_logs,omitempty"`
}

type UserResponse struct {
	ID        uint      `json:"id"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	Avatar    string    `json:"avatar"`
	IsAdmin   bool      `json:"is_admin"`
	CreatedAt time.Time `json:"created_at"`
}

func (u *User) ToResponse() UserResponse {
	return UserResponse{
		ID:        u.ID,
		Email:     u.Email,
		Name:      u.Name,
		Avatar:    u.Avatar,
		IsAdmin:   u.IsAdmin,
		CreatedAt: u.CreatedAt,
	}
}
