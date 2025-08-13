package handlers

import (
	"net/http"
	"strconv"

	"digital-wellbeing-backend/internal/models"
	"digital-wellbeing-backend/internal/services"

	"github.com/gin-gonic/gin"
)

type CheckInHandler struct {
	checkInService *services.CheckInService
}

func NewCheckInHandler(checkInService *services.CheckInService) *CheckInHandler {
	return &CheckInHandler{checkInService: checkInService}
}

func (h *CheckInHandler) CreateCheckIn(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req models.CheckInRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	checkIn, err := h.checkInService.CreateCheckIn(userID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create check-in"})
		return
	}

	c.JSON(http.StatusCreated, checkIn)
}

func (h *CheckInHandler) GetCheckIns(c *gin.Context) {
	userID := c.GetUint("user_id")

	limitStr := c.DefaultQuery("limit", "30")
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		limit = 30
	}

	checkIns, err := h.checkInService.GetUserCheckIns(userID, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch check-ins"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"check_ins": checkIns})
}

func (h *CheckInHandler) GetTodayCheckIn(c *gin.Context) {
	userID := c.GetUint("user_id")

	checkIn, err := h.checkInService.GetTodayCheckIn(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "No check-in found for today"})
		return
	}

	c.JSON(http.StatusOK, checkIn)
}

func (h *CheckInHandler) GetMoodStats(c *gin.Context) {
	userID := c.GetUint("user_id")

	daysStr := c.DefaultQuery("days", "30")
	days, err := strconv.Atoi(daysStr)
	if err != nil {
		days = 30
	}

	stats, err := h.checkInService.GetMoodStats(userID, days)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch mood stats"})
		return
	}

	c.JSON(http.StatusOK, stats)
}
