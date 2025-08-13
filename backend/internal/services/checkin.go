package services

import (
	"time"

	"digital-wellbeing-backend/internal/models"

	"gorm.io/gorm"
)

type CheckInService struct {
	db *gorm.DB
}

func NewCheckInService(db *gorm.DB) *CheckInService {
	return &CheckInService{db: db}
}

func (s *CheckInService) CreateCheckIn(userID uint, req models.CheckInRequest) (*models.CheckIn, error) {
	today := time.Now().Format("2006-01-02")

	// Check if user already checked in today
	var existingCheckIn models.CheckIn
	if err := s.db.Where("user_id = ? AND DATE(date) = ?", userID, today).First(&existingCheckIn).Error; err == nil {
		// Update existing check-in
		existingCheckIn.Mood = req.Mood
		existingCheckIn.ScreenTime = req.ScreenTime
		existingCheckIn.Notes = req.Notes

		if err := s.db.Save(&existingCheckIn).Error; err != nil {
			return nil, err
		}
		return &existingCheckIn, nil
	}

	// Create new check-in
	checkIn := models.CheckIn{
		UserID:     userID,
		Date:       time.Now(),
		Mood:       req.Mood,
		ScreenTime: req.ScreenTime,
		Notes:      req.Notes,
	}

	if err := s.db.Create(&checkIn).Error; err != nil {
		return nil, err
	}

	return &checkIn, nil
}

func (s *CheckInService) GetUserCheckIns(userID uint, limit int) ([]models.CheckIn, error) {
	var checkIns []models.CheckIn
	query := s.db.Where("user_id = ?", userID).Order("date DESC")

	if limit > 0 {
		query = query.Limit(limit)
	}

	if err := query.Find(&checkIns).Error; err != nil {
		return nil, err
	}

	return checkIns, nil
}

func (s *CheckInService) GetTodayCheckIn(userID uint) (*models.CheckIn, error) {
	var checkIn models.CheckIn
	today := time.Now().Format("2006-01-02")

	if err := s.db.Where("user_id = ? AND DATE(date) = ?", userID, today).First(&checkIn).Error; err != nil {
		return nil, err
	}

	return &checkIn, nil
}

func (s *CheckInService) GetMoodStats(userID uint, days int) (*models.MoodStats, error) {
	var checkIns []models.CheckIn

	startDate := time.Now().AddDate(0, 0, -days)
	if err := s.db.Where("user_id = ? AND date >= ?", userID, startDate).Find(&checkIns).Error; err != nil {
		return nil, err
	}

	if len(checkIns) == 0 {
		return &models.MoodStats{
			MoodDistribution: make(map[int]int),
		}, nil
	}

	// Calculate average mood
	var totalMood int
	moodDistribution := make(map[int]int)

	for _, checkIn := range checkIns {
		totalMood += checkIn.Mood
		moodDistribution[checkIn.Mood]++
	}

	averageMood := float64(totalMood) / float64(len(checkIns))

	// Calculate streak (consecutive days with check-ins)
	streak := s.calculateStreak(userID)

	return &models.MoodStats{
		AverageMood:      averageMood,
		MoodDistribution: moodDistribution,
		Streak:           streak,
		TotalCheckIns:    len(checkIns),
	}, nil
}

func (s *CheckInService) calculateStreak(userID uint) int {
	var checkIns []models.CheckIn

	// Get check-ins ordered by date descending
	if err := s.db.Where("user_id = ?", userID).Order("date DESC").Find(&checkIns).Error; err != nil {
		return 0
	}

	if len(checkIns) == 0 {
		return 0
	}

	streak := 0
	today := time.Now()

	for i, checkIn := range checkIns {
		expectedDate := today.AddDate(0, 0, -i)
		checkInDate := checkIn.Date.Format("2006-01-02")
		expectedDateStr := expectedDate.Format("2006-01-02")

		if checkInDate == expectedDateStr {
			streak++
		} else {
			break
		}
	}

	return streak
}
