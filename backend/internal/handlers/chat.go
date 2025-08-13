package handlers

import (
	"net/http"
	"strconv"

	"digital-wellbeing-backend/internal/models"
	"digital-wellbeing-backend/internal/services"

	"github.com/gin-gonic/gin"
)

type ChatHandler struct {
	chatService *services.ChatService
}

func NewChatHandler(chatService *services.ChatService) *ChatHandler {
	return &ChatHandler{chatService: chatService}
}

func (h *ChatHandler) SendMessage(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req models.ChatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := h.chatService.ProcessMessage(userID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process message"})
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *ChatHandler) GetHistory(c *gin.Context) {
	userID := c.GetUint("user_id")

	limitStr := c.DefaultQuery("limit", "50")
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		limit = 50
	}

	history, err := h.chatService.GetChatHistory(userID, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch chat history"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"history": history})
}
