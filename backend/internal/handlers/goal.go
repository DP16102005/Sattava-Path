package handlers

import (
	"net/http"
	"strconv"

	"digital-wellbeing-backend/internal/models"
	"digital-wellbeing-backend/internal/services"

	"github.com/gin-gonic/gin"
)

type GoalHandler struct {
	goalService *services.GoalService
}

func NewGoalHandler(goalService *services.GoalService) *GoalHandler {
	return &GoalHandler{goalService: goalService}
}

func (h *GoalHandler) CreateGoal(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req models.GoalRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	goal, err := h.goalService.CreateGoal(userID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create goal"})
		return
	}

	c.JSON(http.StatusCreated, goal)
}

func (h *GoalHandler) GetGoals(c *gin.Context) {
	userID := c.GetUint("user_id")

	goals, err := h.goalService.GetUserGoals(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch goals"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"goals": goals})
}

func (h *GoalHandler) UpdateGoal(c *gin.Context) {
	userID := c.GetUint("user_id")
	goalIDStr := c.Param("id")

	goalID, err := strconv.ParseUint(goalIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid goal ID"})
		return
	}

	var req models.GoalRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	goal, err := h.goalService.UpdateGoal(userID, uint(goalID), req)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Goal not found"})
		return
	}

	c.JSON(http.StatusOK, goal)
}

func (h *GoalHandler) DeleteGoal(c *gin.Context) {
	userID := c.GetUint("user_id")
	goalIDStr := c.Param("id")

	goalID, err := strconv.ParseUint(goalIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid goal ID"})
		return
	}

	if err := h.goalService.DeleteGoal(userID, uint(goalID)); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Goal not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Goal deleted successfully"})
}

func (h *GoalHandler) GetProgress(c *gin.Context) {
	userID := c.GetUint("user_id")

	progress, err := h.goalService.GetGoalProgress(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch progress"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"progress": progress})
}
