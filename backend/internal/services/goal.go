package services

import (
	"time"

	"digital-wellbeing-backend/internal/models"

	"gorm.io/gorm"
)

type GoalService struct {
	db *gorm.DB
}

func NewGoalService(db *gorm.DB) *GoalService {
	return &GoalService{db: db}
}

func (s *GoalService) CreateGoal(userID uint, req models.GoalRequest) (*models.Goal, error) {
	goal := models.Goal{
		UserID:     userID,
		AppName:    req.AppName,
		AppIcon:    req.AppIcon,
		DailyLimit: req.DailyLimit,
		IsEnabled:  req.IsEnabled,
	}

	if err := s.db.Create(&goal).Error; err != nil {
		return nil, err
	}

	return &goal, nil
}

func (s *GoalService) GetUserGoals(userID uint) ([]models.Goal, error) {
	var goals []models.Goal
	if err := s.db.Where("user_id = ?", userID).Find(&goals).Error; err != nil {
		return nil, err
	}
	return goals, nil
}

func (s *GoalService) UpdateGoal(userID uint, goalID uint, req models.GoalRequest) (*models.Goal, error) {
	var goal models.Goal
	if err := s.db.Where("id = ? AND user_id = ?", goalID, userID).First(&goal).Error; err != nil {
		return nil, err
	}

	goal.AppName = req.AppName
	goal.AppIcon = req.AppIcon
	goal.DailyLimit = req.DailyLimit
	goal.IsEnabled = req.IsEnabled

	if err := s.db.Save(&goal).Error; err != nil {
		return nil, err
	}

	return &goal, nil
}

func (s *GoalService) DeleteGoal(userID uint, goalID uint) error {
	return s.db.Where("id = ? AND user_id = ?", goalID, userID).Delete(&models.Goal{}).Error
}

func (s *GoalService) GetGoalProgress(userID uint) ([]models.GoalProgress, error) {
	goals, err := s.GetUserGoals(userID)
	if err != nil {
		return nil, err
	}

	var progress []models.GoalProgress
	today := time.Now().Format("2006-01-02")

	for _, goal := range goals {
		if !goal.IsEnabled {
			continue
		}

		// Get today's usage for this app
		var totalUsage int
		s.db.Model(&models.UsageLog{}).
			Where("user_id = ? AND app_name = ? AND DATE(date) = ?", userID, goal.AppName, today).
			Select("COALESCE(SUM(duration), 0)").
			Scan(&totalUsage)

		remaining := goal.DailyLimit - totalUsage
		if remaining < 0 {
			remaining = 0
		}

		progressPercent := float64(totalUsage) / float64(goal.DailyLimit) * 100
		if progressPercent > 100 {
			progressPercent = 100
		}

		progress = append(progress, models.GoalProgress{
			Goal:            goal,
			UsedToday:       totalUsage,
			RemainingTime:   remaining,
			IsOverLimit:     totalUsage > goal.DailyLimit,
			ProgressPercent: progressPercent,
		})
	}

	return progress, nil
}
