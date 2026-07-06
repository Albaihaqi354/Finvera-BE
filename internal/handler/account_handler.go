package handler

import (
	"finvera-be/internal/dto"
	"finvera-be/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AccountHandler struct {
	accountService service.AccountService
}

func NewAccountHandler(accountService service.AccountService) *AccountHandler {
	return &AccountHandler{accountService: accountService}
}

// @Summary Create a new account
// @Description Create a new account for the authenticated user
// @Tags Accounts
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.CreateAccountRequest true "Create Account Request"
// @Success 201 {object} dto.Response
// @Failure 400 {object} dto.Response
// @Failure 401 {object} dto.Response
// @Router /accounts [post]
func (h *AccountHandler) CreateAccount(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse("Invalid user ID in token"))
		return
	}

	var req dto.CreateAccountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse(err.Error()))
		return
	}

	account, err := h.accountService.CreateAccount(userID, req)
	if err != nil {
		handleServiceError(c, err)
		return
	}

	c.JSON(http.StatusCreated, dto.SuccessResponse("Account created successfully", account))
}

// @Summary Get all accounts
// @Description Get all accounts for the authenticated user
// @Tags Accounts
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number"
// @Param limit query int false "Items per page"
// @Success 200 {object} dto.Response
// @Failure 401 {object} dto.Response
// @Router /accounts [get]
func (h *AccountHandler) GetAccounts(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse("Invalid user ID in token"))
		return
	}

	page, limit := dto.GetPaginationParams(c)

	accounts, total, err := h.accountService.GetAccounts(userID, page, limit)
	if err != nil {
		handleServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, dto.PaginatedResponse("Accounts retrieved successfully", accounts, page, limit, total))
}

// @Summary Get account by ID
// @Description Get account details by account ID
// @Tags Accounts
// @Produce json
// @Security BearerAuth
// @Param id path string true "Account ID"
// @Success 200 {object} dto.Response
// @Failure 401 {object} dto.Response
// @Failure 404 {object} dto.Response
// @Router /accounts/{id} [get]
func (h *AccountHandler) GetAccountByID(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse("Invalid user ID in token"))
		return
	}

	accountIDStr := c.Param("id")
	accountID, err := uuid.Parse(accountIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse("Invalid account ID format"))
		return
	}

	account, err := h.accountService.GetAccountByID(userID, accountID)
	if err != nil {
		handleServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse("Account retrieved successfully", account))
}

// @Summary Update account
// @Description Update account details
// @Tags Accounts
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Account ID"
// @Param request body dto.UpdateAccountRequest true "Update Account Request"
// @Success 200 {object} dto.Response
// @Failure 400 {object} dto.Response
// @Failure 401 {object} dto.Response
// @Failure 404 {object} dto.Response
// @Router /accounts/{id} [put]
func (h *AccountHandler) UpdateAccount(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse("Invalid user ID in token"))
		return
	}

	accountIDStr := c.Param("id")
	accountID, err := uuid.Parse(accountIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse("Invalid account ID format"))
		return
	}

	var req dto.UpdateAccountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse(err.Error()))
		return
	}

	account, err := h.accountService.UpdateAccount(userID, accountID, req)
	if err != nil {
		handleServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse("Account updated successfully", account))
}

// @Summary Delete account
// @Description Delete an account
// @Tags Accounts
// @Produce json
// @Security BearerAuth
// @Param id path string true "Account ID"
// @Success 200 {object} dto.Response
// @Failure 401 {object} dto.Response
// @Failure 404 {object} dto.Response
// @Router /accounts/{id} [delete]
func (h *AccountHandler) DeleteAccount(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse("Invalid user ID in token"))
		return
	}

	accountIDStr := c.Param("id")
	accountID, err := uuid.Parse(accountIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse("Invalid account ID format"))
		return
	}

	if err := h.accountService.DeleteAccount(userID, accountID); err != nil {
		handleServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse("Account deleted successfully", nil))
}
