package handler

import (
	"finvera-be/internal/dto"
	"finvera-be/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CategoryHandler struct {
	categoryService service.CategoryService
}

func NewCategoryHandler(categoryService service.CategoryService) *CategoryHandler {
	return &CategoryHandler{categoryService: categoryService}
}

// @Summary Create a new category
// @Description Create a new custom category for the authenticated user
// @Tags Categories
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.CreateCategoryRequest true "Create Category Request"
// @Success 201 {object} dto.Response
// @Failure 400 {object} dto.Response
// @Failure 401 {object} dto.Response
// @Router /categories [post]
func (h *CategoryHandler) CreateCategory(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse("Invalid user ID in token"))
		return
	}

	var req dto.CreateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse(err.Error()))
		return
	}

	category, err := h.categoryService.CreateCategory(userID, req)
	if err != nil {
		handleServiceError(c, err)
		return
	}

	c.JSON(http.StatusCreated, dto.SuccessResponse("Category created successfully", category))
}

// @Summary Get all categories
// @Description Get all default and custom categories for the user
// @Tags Categories
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number"
// @Param limit query int false "Items per page"
// @Success 200 {object} dto.Response
// @Failure 401 {object} dto.Response
// @Router /categories [get]
func (h *CategoryHandler) GetCategories(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse("Invalid user ID in token"))
		return
	}

	page, limit := dto.GetPaginationParams(c)

	categories, total, err := h.categoryService.GetCategories(userID, page, limit)
	if err != nil {
		handleServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, dto.PaginatedResponse("Categories retrieved successfully", categories, page, limit, total))
}

// @Summary Get category by ID
// @Description Get category details by category ID
// @Tags Categories
// @Produce json
// @Security BearerAuth
// @Param id path string true "Category ID"
// @Success 200 {object} dto.Response
// @Failure 400 {object} dto.Response
// @Failure 401 {object} dto.Response
// @Failure 403 {object} dto.Response
// @Failure 404 {object} dto.Response
// @Router /categories/{id} [get]
func (h *CategoryHandler) GetCategoryByID(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse("Invalid user ID in token"))
		return
	}

	categoryIDStr := c.Param("id")
	categoryID, err := uuid.Parse(categoryIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse("Invalid category ID format"))
		return
	}

	category, err := h.categoryService.GetCategoryByID(userID, categoryID)
	if err != nil {
		handleServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse("Category retrieved successfully", category))
}

// @Summary Update category
// @Description Update custom category details
// @Tags Categories
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Category ID"
// @Param request body dto.UpdateCategoryRequest true "Update Category Request"
// @Success 200 {object} dto.Response
// @Failure 400 {object} dto.Response
// @Failure 401 {object} dto.Response
// @Failure 403 {object} dto.Response
// @Failure 404 {object} dto.Response
// @Router /categories/{id} [put]
func (h *CategoryHandler) UpdateCategory(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse("Invalid user ID in token"))
		return
	}

	categoryIDStr := c.Param("id")
	categoryID, err := uuid.Parse(categoryIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse("Invalid category ID format"))
		return
	}

	var req dto.UpdateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse(err.Error()))
		return
	}

	category, err := h.categoryService.UpdateCategory(userID, categoryID, req)
	if err != nil {
		handleServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse("Category updated successfully", category))
}

// @Summary Delete category
// @Description Delete a custom category
// @Tags Categories
// @Produce json
// @Security BearerAuth
// @Param id path string true "Category ID"
// @Success 200 {object} dto.Response
// @Failure 400 {object} dto.Response
// @Failure 401 {object} dto.Response
// @Failure 403 {object} dto.Response
// @Failure 404 {object} dto.Response
// @Router /categories/{id} [delete]
func (h *CategoryHandler) DeleteCategory(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse("Invalid user ID in token"))
		return
	}

	categoryIDStr := c.Param("id")
	categoryID, err := uuid.Parse(categoryIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse("Invalid category ID format"))
		return
	}

	if err := h.categoryService.DeleteCategory(userID, categoryID); err != nil {
		handleServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse("Category deleted successfully", nil))
}

// @Summary Get preset categories
// @Description Get all system preset categories (grouped with children)
// @Tags Categories
// @Produce json
// @Success 200 {object} dto.Response
// @Router /preset-categories [get]
func (h *CategoryHandler) GetPresetCategories(c *gin.Context) {
	categories, err := h.categoryService.GetPresetCategories()
	if err != nil {
		handleServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse("Preset categories retrieved successfully", categories))
}
