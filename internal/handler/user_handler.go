package handler

import (
	"finvera-be/internal/dto"
	"finvera-be/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UserHandler struct {
	userService service.UserService
}

func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

// @Summary Get user profile
// @Description Get profile of the currently authenticated user
// @Tags Users
// @Produce json
// @Security BearerAuth
// @Success 200 {object} dto.Response
// @Failure 401 {object} dto.Response
// @Failure 500 {object} dto.Response
// @Router /users/me [get]
func (h *UserHandler) GetProfile(c *gin.Context) {
	userIdStr, _ := c.Get("userId")
	userID, err := uuid.Parse(userIdStr.(string))
	if err != nil {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse("Invalid user ID in token"))
		return
	}

	user, err := h.userService.GetProfile(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse("Profile retrieved successfully", user))
}

// @Summary Update user profile
// @Description Update username or email of the currently authenticated user
// @Tags Users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.UpdateUserRequest true "Update User Request"
// @Success 200 {object} dto.Response
// @Failure 400 {object} dto.Response
// @Failure 401 {object} dto.Response
// @Failure 500 {object} dto.Response
// @Router /users/me [put]
func (h *UserHandler) UpdateProfile(c *gin.Context) {
	userIdStr, _ := c.Get("userId")
	userID, err := uuid.Parse(userIdStr.(string))
	if err != nil {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse("Invalid user ID in token"))
		return
	}

	var req dto.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse(err.Error()))
		return
	}

	user, err := h.userService.UpdateProfile(userID, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse("Profile updated successfully", user))
}
