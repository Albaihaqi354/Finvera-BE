package handler

import (
	"finvera-be/internal/dto"
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
// @Param request body dto.CreateScheduledRequest true "Create Scheduled Transaction Request"
// @Success 201 {object} dto.Response
// @Failure 400 {object} dto.Response
// @Failure 401 {object} dto.Response
// @Failure 500 {object} dto.Response
// @Router /scheduled [post]
func (h *ScheduledTransactionHandler) CreateScheduled(c *gin.Context) {
	userIdStr, _ := c.Get("userId")
	userID, err := uuid.Parse(userIdStr.(string))
	if err != nil {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse("Invalid user ID in token"))
		return
	}

	var req dto.CreateScheduledRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse(err.Error()))
		return
	}

	scheduled, err := h.scheduledService.CreateScheduled(userID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse(err.Error()))
		return
	}

	c.JSON(http.StatusCreated, dto.SuccessResponse("Scheduled transaction created successfully", scheduled))
}

// @Summary Get all scheduled transactions
// @Description Get all scheduled transactions for the authenticated user
// @Tags ScheduledTransactions
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number"
// @Param limit query int false "Items per page"
// @Success 200 {object} dto.Response
// @Failure 401 {object} dto.Response
// @Router /scheduled [get]
func (h *ScheduledTransactionHandler) GetScheduleds(c *gin.Context) {
	userIdStr, _ := c.Get("userId")
	userID, err := uuid.Parse(userIdStr.(string))
	if err != nil {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse("Invalid user ID in token"))
		return
	}

	page, limit := dto.GetPaginationParams(c)

	scheduleds, total, err := h.scheduledService.GetScheduleds(userID, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, dto.PaginatedResponse("Scheduled transactions retrieved successfully", scheduleds, page, limit, total))
}

// @Summary Get scheduled transaction by ID
// @Description Get scheduled transaction details by ID
// @Tags ScheduledTransactions
// @Produce json
// @Security BearerAuth
// @Param id path string true "Scheduled Transaction ID"
// @Success 200 {object} dto.Response
// @Failure 400 {object} dto.Response
// @Failure 401 {object} dto.Response
// @Failure 404 {object} dto.Response
// @Router /scheduled/{id} [get]
func (h *ScheduledTransactionHandler) GetScheduledByID(c *gin.Context) {
	userIdStr, _ := c.Get("userId")
	userID, err := uuid.Parse(userIdStr.(string))
	if err != nil {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse("Invalid user ID in token"))
		return
	}

	scheduledIDStr := c.Param("id")
	scheduledID, err := uuid.Parse(scheduledIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse("Invalid scheduled transaction ID format"))
		return
	}

	scheduled, err := h.scheduledService.GetScheduledByID(userID, scheduledID)
	if err != nil {
		c.JSON(http.StatusNotFound, dto.ErrorResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse("Scheduled transaction retrieved successfully", scheduled))
}

// @Summary Update scheduled transaction
// @Description Update scheduled transaction details
// @Tags ScheduledTransactions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Scheduled Transaction ID"
// @Param request body dto.UpdateScheduledRequest true "Update Scheduled Transaction Request"
// @Success 200 {object} dto.Response
// @Failure 400 {object} dto.Response
// @Failure 401 {object} dto.Response
// @Failure 500 {object} dto.Response
// @Router /scheduled/{id} [put]
func (h *ScheduledTransactionHandler) UpdateScheduled(c *gin.Context) {
	userIdStr, _ := c.Get("userId")
	userID, err := uuid.Parse(userIdStr.(string))
	if err != nil {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse("Invalid user ID in token"))
		return
	}

	scheduledIDStr := c.Param("id")
	scheduledID, err := uuid.Parse(scheduledIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse("Invalid scheduled transaction ID format"))
		return
	}

	var req dto.UpdateScheduledRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse(err.Error()))
		return
	}

	scheduled, err := h.scheduledService.UpdateScheduled(userID, scheduledID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse("Scheduled transaction updated successfully", scheduled))
}

// @Summary Delete scheduled transaction
// @Description Delete a scheduled transaction
// @Tags ScheduledTransactions
// @Produce json
// @Security BearerAuth
// @Param id path string true "Scheduled Transaction ID"
// @Success 200 {object} dto.Response
// @Failure 400 {object} dto.Response
// @Failure 401 {object} dto.Response
// @Failure 500 {object} dto.Response
// @Router /scheduled/{id} [delete]
func (h *ScheduledTransactionHandler) DeleteScheduled(c *gin.Context) {
	userIdStr, _ := c.Get("userId")
	userID, err := uuid.Parse(userIdStr.(string))
	if err != nil {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse("Invalid user ID in token"))
		return
	}

	scheduledIDStr := c.Param("id")
	scheduledID, err := uuid.Parse(scheduledIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse("Invalid scheduled transaction ID format"))
		return
	}

	if err := h.scheduledService.DeleteScheduled(userID, scheduledID); err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse("Scheduled transaction deleted successfully", nil))
}
