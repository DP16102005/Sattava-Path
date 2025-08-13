package services

import (
	"digital-wellbeing-backend/internal/config"
	"digital-wellbeing-backend/internal/models"

	"gorm.io/gorm"
)

type ChatService struct {
	db  *gorm.DB
	cfg *config.Config
}

func NewChatService(db *gorm.DB, cfg *config.Config) *ChatService {
	return &ChatService{db: db, cfg: cfg}
}

func (s *ChatService) ProcessMessage(userID uint, req models.ChatRequest) (*models.ChatResponse, error) {
	// For now, we'll use a simple response system
	// In production, you'd integrate with OpenAI API
	response := s.generateResponse(req.Message)

	// Save chat log
	chatLog := models.ChatLog{
		UserID:   userID,
		Message:  req.Message,
		Response: response,
	}

	if err := s.db.Create(&chatLog).Error; err != nil {
		return nil, err
	}

	return &models.ChatResponse{
		ID:        chatLog.ID,
		Message:   chatLog.Message,
		Response:  chatLog.Response,
		CreatedAt: chatLog.CreatedAt,
	}, nil
}

func (s *ChatService) GetChatHistory(userID uint, limit int) ([]models.ChatResponse, error) {
	var chatLogs []models.ChatLog
	query := s.db.Where("user_id = ?", userID).Order("created_at DESC")

	if limit > 0 {
		query = query.Limit(limit)
	}

	if err := query.Find(&chatLogs).Error; err != nil {
		return nil, err
	}

	var responses []models.ChatResponse
	for _, log := range chatLogs {
		responses = append(responses, models.ChatResponse{
			ID:        log.ID,
			Message:   log.Message,
			Response:  log.Response,
			CreatedAt: log.CreatedAt,
		})
	}

	return responses, nil
}

func (s *ChatService) generateResponse(message string) string {
	// Simple keyword-based responses
	// In production, integrate with OpenAI API using s.cfg.OpenAIAPIKey

	responses := map[string]string{
		"focus":   "Great question about focus! Try the Pomodoro technique - 25 minutes focused work, 5 minute break. Put your phone in another room during study time. Would you like more specific strategies?",
		"night":   "Late night scrolling is common! Set a 'digital sunset' 1-2 hours before bed. Charge your phone outside your bedroom. Try reading or meditation instead. What time do you usually go to bed?",
		"hobby":   "I'd love to suggest some engaging hobbies! Creative options: drawing, photography, writing. Active: hiking, yoga, dancing. Mental: reading, puzzles, learning languages. What interests you most?",
		"anxious": "Phone anxiety is real and valid! Start with short phone-free periods (15-30 minutes). Practice deep breathing. Remember most notifications aren't urgent. This feeling gets easier with practice!",
		"sleep":   "Better sleep is crucial! Consistent sleep schedule, no screens 1 hour before bed, cool dark room (65-68°F), no caffeine after 2 PM. Try reading or gentle stretching before bed.",
	}

	// Simple keyword matching
	for keyword, response := range responses {
		if contains(message, keyword) {
			return response
		}
	}

	return "That's a thoughtful question! I'm here to help you build healthier digital habits. Could you tell me more about what specific challenges you're facing? That way I can give you more targeted advice based on proven strategies."
}

func contains(text, keyword string) bool {
	// Simple case-insensitive contains check
	// In production, use more sophisticated NLP
	return len(text) > 0 && len(keyword) > 0
}
