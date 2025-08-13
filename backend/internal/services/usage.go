package services

import (
	"time"

	"digital-wellbeing-backend/internal/models"

	"gorm.io/gorm"
)

type UsageService struct {
	db *gorm.DB
}

func NewUsageService(db *gorm.DB) *UsageService {
	return &UsageService{db: db}
}

func (s *UsageService) LogUsage(userID uint, req models.UsageRequest) (*models.UsageLog, error) {
	today := time.Now().Format("2006-01-02")

	// Check if there's already a log for this app today
	var existingLog models.UsageLog
	if err := s.db.Where("user_id = ? AND app_name = ? AND DATE(date) = ?", userID, req.AppName, today).First(&existingLog).Error; err == nil {
		// Update existing log
		existingLog.Duration += req.Duration
		if err := s.db.Save(&existingLog).Error; err != nil {
			return nil, err
		}
		return &existingLog, nil
	}

	// Create new log
	usageLog := models.UsageLog{
		UserID:   userID,
		AppName:  req.AppName,
		Date:     time.Now(),
		Duration: req.Duration,
	}

	if err := s.db.Create(&usageLog).Error; err != nil {
		return nil, err
	}

	return &usageLog, nil
}

func (s *UsageService) GetDailyUsage(userID uint, date time.Time) ([]models.DailyUsage, error) {
	var usageLogs []models.UsageLog
	dateStr := date.Format("2006-01-02")

	if err := s.db.Where("user_id = ? AND DATE(date) = ?", userID, dateStr).Find(&usageLogs).Error; err != nil {
		return nil, err
	}

	var dailyUsage []models.DailyUsage
	for _, log := range usageLogs {
		dailyUsage = append(dailyUsage, models.DailyUsage{
			Date:     log.Date,
			AppName:  log.AppName,
			Duration: log.Duration,
		})
	}

	return dailyUsage, nil
}

func (s *UsageService) GetWeeklyUsage(userID uint) ([]models.DailyUsage, error) {
	var usageLogs []models.UsageLog
	weekAgo := time.Now().AddDate(0, 0, -7)

	if err := s.db.Where("user_id = ? AND date >= ?", userID, weekAgo).Find(&usageLogs).Error; err != nil {
		return nil, err
	}

	var weeklyUsage []models.DailyUsage
	for _, log := range usageLogs {
		weeklyUsage = append(weeklyUsage, models.DailyUsage{
			Date:     log.Date,
			AppName:  log.AppName,
			Duration: log.Duration,
		})
	}

	return weeklyUsage, nil
}

func (s *UsageService) GetUsageStats(userID uint, days int) (*models.UsageStats, error) {
	var usageLogs []models.UsageLog
	startDate := time.Now().AddDate(0, 0, -days)

	if err := s.db.Where("user_id = ? AND date >= ?", userID, startDate).Find(&usageLogs).Error; err != nil {
		return nil, err
	}

	if len(usageLogs) == 0 {
		return &models.UsageStats{
			AppBreakdown: make(map[string]int),
		}, nil
	}

	// Calculate stats
	totalScreenTime := 0
	appBreakdown := make(map[string]int)
	mostUsedApp := ""
	maxUsage := 0

	for _, log := range usageLogs {
		totalScreenTime += log.Duration
		appBreakdown[log.AppName] += log.Duration

		if appBreakdown[log.AppName] > maxUsage {
			maxUsage = appBreakdown[log.AppName]
			mostUsedApp = log.AppName
		}
	}

	averageDaily := totalScreenTime / days
	if days == 0 {
		averageDaily = 0
	}

	// Get weekly trend
	weeklyTrend, _ := s.GetWeeklyUsage(userID)

	return &models.UsageStats{
		TotalScreenTime: totalScreenTime,
		AverageDaily:    averageDaily,
		MostUsedApp:     mostUsedApp,
		AppBreakdown:    appBreakdown,
		WeeklyTrend:     weeklyTrend,
	}, nil
}
