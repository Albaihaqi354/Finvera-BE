package handler

import (
	"net/http"

	"finvera-be/internal/dto"
	"finvera-be/internal/repository"
	"finvera-be/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// allowedTransactionTypes is the whitelist for the type query filter.
var allowedTransactionTypes = map[string]bool{
	"income":   true,
	"expense":  true,
	"transfer": true,
}

type TransactionHandler struct {
	transactionService service.TransactionService
}

func NewTransactionHandler(transactionService service.TransactionService) *TransactionHandler {
	return &TransactionHandler{transactionService: transactionService}
}

// getUserID extracts and parses the userId set by AuthMiddleware.
func getUserID(c *gin.Context) (uuid.UUID, bool) {
	v, exists := c.Get("userId")
	if !exists {
		return uuid.Nil, false
	}
	id, err := uuid.Parse(v.(string))
	if err != nil {
		return uuid.Nil, false
	}
	return id, true
}

// @Summary Create a new transaction
// @Tags Transactions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.CreateTransactionRequest true "Create Transaction Request"
// @Success 201 {object} dto.Response
// @Failure 400 {object} dto.Response
// @Failure 401 {object} dto.Response
// @Failure 500 {object} dto.Response
// @Router /transactions [post]
func (h *TransactionHandler) CreateTransaction(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse("Invalid or missing user ID in token"))
		return
	}

	var req dto.CreateTransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse(err.Error()))
		return
	}

	transaction, err := h.transactionService.CreateTransaction(userID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse(err.Error()))
		return
	}

	c.JSON(http.StatusCreated, dto.SuccessResponse("Transaction created successfully", transaction))
}

// @Summary Get all transactions
// @Tags Transactions
// @Produce json
// @Security BearerAuth
// @Param page     query int    false "Page number (default: 1)"
// @Param limit    query int    false "Items per page (default: 10, max: 100)"
// @Param startDate query string false "Start date (RFC3339)"
// @Param endDate   query string false "End date (RFC3339)"
// @Param type      query string false "Transaction type: income|expense|transfer"
// @Param accountId query string false "Account UUID"
// @Param search    query string false "Search by note"
// @Success 200 {object} dto.Response
// @Failure 401 {object} dto.Response
// @Router /transactions [get]
func (h *TransactionHandler) GetTransactions(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse("Invalid or missing user ID in token"))
		return
	}

	page, limit := dto.GetPaginationParams(c)

	// Validate & sanitize type filter — whitelist only
	typeFilter := c.Query("type")
	if typeFilter != "" && !allowedTransactionTypes[typeFilter] {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse("Invalid type filter. Must be one of: income, expense, transfer"))
		return
	}

	// Validate accountId is a valid UUID if provided
	accountIDStr := c.Query("accountId")
	if accountIDStr != "" {
		if _, err := uuid.Parse(accountIDStr); err != nil {
			c.JSON(http.StatusBadRequest, dto.ErrorResponse("Invalid accountId format"))
			return
		}
	}

	filter := repository.TransactionFilter{
		StartDate: c.Query("startDate"),
		EndDate:   c.Query("endDate"),
		Type:      typeFilter,
		AccountID: accountIDStr,
		Search:    c.Query("search"),
	}

	transactions, total, err := h.transactionService.GetTransactions(userID, page, limit, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, dto.PaginatedResponse("Transactions retrieved successfully", transactions, page, limit, total))
}

// @Summary Get transaction by ID
// @Tags Transactions
// @Produce json
// @Security BearerAuth
// @Param id path string true "Transaction UUID"
// @Success 200 {object} dto.Response
// @Failure 400 {object} dto.Response
// @Failure 401 {object} dto.Response
// @Failure 404 {object} dto.Response
// @Router /transactions/{id} [get]
func (h *TransactionHandler) GetTransactionByID(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse("Invalid or missing user ID in token"))
		return
	}

	transactionID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse("Invalid transaction ID format"))
		return
	}

	transaction, err := h.transactionService.GetTransactionByID(userID, transactionID)
	if err != nil {
		c.JSON(http.StatusNotFound, dto.ErrorResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse("Transaction retrieved successfully", transaction))
}

// @Summary Update transaction
// @Tags Transactions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Transaction UUID"
// @Param request body dto.TransactionRequest true "Update Transaction Request"
// @Success 200 {object} dto.Response
// @Failure 400 {object} dto.Response
// @Failure 401 {object} dto.Response
// @Failure 500 {object} dto.Response
// @Router /transactions/{id} [put]
func (h *TransactionHandler) UpdateTransaction(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse("Invalid or missing user ID in token"))
		return
	}

	transactionID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse("Invalid transaction ID format"))
		return
	}

	var req dto.TransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse(err.Error()))
		return
	}

	transaction, err := h.transactionService.UpdateTransaction(userID, transactionID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse("Transaction updated successfully", transaction))
}

// @Summary Delete transaction
// @Tags Transactions
// @Produce json
// @Security BearerAuth
// @Param id path string true "Transaction UUID"
// @Success 200 {object} dto.Response
// @Failure 400 {object} dto.Response
// @Failure 401 {object} dto.Response
// @Failure 500 {object} dto.Response
// @Router /transactions/{id} [delete]
func (h *TransactionHandler) DeleteTransaction(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse("Invalid or missing user ID in token"))
		return
	}

	transactionID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse("Invalid transaction ID format"))
		return
	}

	if err := h.transactionService.DeleteTransaction(userID, transactionID); err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse("Transaction deleted successfully", nil))
}
