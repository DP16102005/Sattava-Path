package handlers

import (
	"net/http"
	"strconv"
	"time"

	"digital-wellbeing-backend/internal/models"
	"digital-wellbeing-backend/internal/services"

	"github.com/gin-gonic/gin"
)

type UsageHandler struct {
	usageService *services.UsageService
}

func NewUsageHandler(usageService *services.UsageService) *UsageHandler {
	return &UsageHandler{usageService: usageService}
}

func (h *UsageHandler) LogUsage(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req models.UsageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	usageLog, err := h.usageService.LogUsage(userID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to log usage"})
		return
	}

	c.JSON(http.StatusCreated, usageLog)
}

func (h *UsageHandler) GetDailyUsage(c *gin.Context) {
	userID := c.GetUint("user_id")

	dateStr := c.DefaultQuery("date", time.Now().Format("2006-01-02"))
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format. Use YYYY-MM-DD"})
		return
	}

	usage, err := h.usageService.GetDailyUsage(userID, date)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch daily usage"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"usage": usage})
}

func (h *UsageHandler) GetWeeklyUsage(c *gin.Context) {
	userID := c.GetUint("user_id")

	usage, err := h.usageService.GetWeeklyUsage(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch weekly usage"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"usage": usage})
}

func (h *UsageHandler) GetStats(c *gin.Context) {
	userID := c.GetUint("user_id")

	daysStr := c.DefaultQuery("days", "30")
	days, err := strconv.Atoi(daysStr)
	if err != nil {
		days = 30
	}

	stats, err := h.usageService.GetUsageStats(userID, days)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch usage stats"})
		return
	}

	c.JSON(http.StatusOK, stats)
}
