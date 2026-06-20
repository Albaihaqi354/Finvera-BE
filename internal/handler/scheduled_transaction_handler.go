package handler

import (
	"finvera-be/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ScheduledTransactionHandler struct {
	scheduledService service.ScheduledTransactionService
}

func NewScheduledTransactionHandler(scheduledService service.ScheduledTransactionService) *ScheduledTransactionHandler {
	return &ScheduledTransactionHandler{scheduledService: scheduledService}
}

// @Summary Create a new scheduled transaction
// @Description Create a new scheduled transaction for the authenticated user
// @Tags ScheduledTransactions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body service.CreateScheduledRequest true "Create Scheduled Transaction Request"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/scheduled [post]
func (h *ScheduledTransactionHandler) CreateScheduled(c *gin.Context) {
	userIdStr, _ := c.Get("userId")
	userID, err := uuid.Parse(userIdStr.(string))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID in token"})
		return
	}

	var req service.CreateScheduledRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	scheduled, err := h.scheduledService.CreateScheduled(userID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Scheduled transaction created successfully", "data": scheduled})
}

// @Summary Get all scheduled transactions
// @Description Get all scheduled transactions for the authenticated user
// @Tags ScheduledTransactions
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /api/v1/scheduled [get]
func (h *ScheduledTransactionHandler) GetScheduleds(c *gin.Context) {
	userIdStr, _ := c.Get("userId")
	userID, err := uuid.Parse(userIdStr.(string))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID in token"})
		return
	}

	scheduleds, err := h.scheduledService.GetScheduleds(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": scheduleds})
}

// @Summary Get scheduled transaction by ID
// @Description Get scheduled transaction details by ID
// @Tags ScheduledTransactions
// @Produce json
// @Security BearerAuth
// @Param id path string true "Scheduled Transaction ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /api/v1/scheduled/{id} [get]
func (h *ScheduledTransactionHandler) GetScheduledByID(c *gin.Context) {
	userIdStr, _ := c.Get("userId")
	userID, err := uuid.Parse(userIdStr.(string))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID in token"})
		return
	}

	scheduledIDStr := c.Param("id")
	scheduledID, err := uuid.Parse(scheduledIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid scheduled transaction ID format"})
		return
	}

	scheduled, err := h.scheduledService.GetScheduledByID(userID, scheduledID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": scheduled})
}

// @Summary Update scheduled transaction
// @Description Update scheduled transaction details
// @Tags ScheduledTransactions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Scheduled Transaction ID"
// @Param request body service.UpdateScheduledRequest true "Update Scheduled Transaction Request"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/scheduled/{id} [put]
func (h *ScheduledTransactionHandler) UpdateScheduled(c *gin.Context) {
	userIdStr, _ := c.Get("userId")
	userID, err := uuid.Parse(userIdStr.(string))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID in token"})
		return
	}

	scheduledIDStr := c.Param("id")
	scheduledID, err := uuid.Parse(scheduledIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid scheduled transaction ID format"})
		return
	}

	var req service.UpdateScheduledRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	scheduled, err := h.scheduledService.UpdateScheduled(userID, scheduledID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Scheduled transaction updated successfully", "data": scheduled})
}

// @Summary Delete scheduled transaction
// @Description Delete a scheduled transaction
// @Tags ScheduledTransactions
// @Produce json
// @Security BearerAuth
// @Param id path string true "Scheduled Transaction ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/scheduled/{id} [delete]
func (h *ScheduledTransactionHandler) DeleteScheduled(c *gin.Context) {
	userIdStr, _ := c.Get("userId")
	userID, err := uuid.Parse(userIdStr.(string))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID in token"})
		return
	}

	scheduledIDStr := c.Param("id")
	scheduledID, err := uuid.Parse(scheduledIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid scheduled transaction ID format"})
		return
	}

	if err := h.scheduledService.DeleteScheduled(userID, scheduledID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Scheduled transaction deleted successfully"})
}
