package handler

import (
	"finvera-be/internal/dto"
	"finvera-be/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type TagHandler struct {
	tagService service.TagService
}

func NewTagHandler(tagService service.TagService) *TagHandler {
	return &TagHandler{tagService: tagService}
}

// @Summary Create a new tag
// @Description Create a new tag for the authenticated user
// @Tags Tags
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.CreateTagRequest true "Create Tag Request"
// @Success 201 {object} dto.Response
// @Failure 400 {object} dto.Response
// @Failure 401 {object} dto.Response
// @Router /tags [post]
func (h *TagHandler) CreateTag(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse("Invalid user ID in token"))
		return
	}

	var req dto.CreateTagRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse(err.Error()))
		return
	}

	tag, err := h.tagService.CreateTag(userID, req)
	if err != nil {
		handleServiceError(c, err)
		return
	}

	c.JSON(http.StatusCreated, dto.SuccessResponse("Tag created successfully", tag))
}

// @Summary Get all tags
// @Description Get all tags for the authenticated user
// @Tags Tags
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number"
// @Param limit query int false "Items per page"
// @Success 200 {object} dto.Response
// @Failure 401 {object} dto.Response
// @Router /tags [get]
func (h *TagHandler) GetTags(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse("Invalid user ID in token"))
		return
	}

	page, limit := dto.GetPaginationParams(c)

	tags, total, err := h.tagService.GetTags(userID, page, limit)
	if err != nil {
		handleServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, dto.PaginatedResponse("Tags retrieved successfully", tags, page, limit, total))
}

// @Summary Get tag by ID
// @Description Get tag details by tag ID
// @Tags Tags
// @Produce json
// @Security BearerAuth
// @Param id path string true "Tag ID"
// @Success 200 {object} dto.Response
// @Failure 400 {object} dto.Response
// @Failure 401 {object} dto.Response
// @Failure 404 {object} dto.Response
// @Router /tags/{id} [get]
func (h *TagHandler) GetTagByID(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse("Invalid user ID in token"))
		return
	}

	tagIDStr := c.Param("id")
	tagID, err := uuid.Parse(tagIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse("Invalid tag ID format"))
		return
	}

	tag, err := h.tagService.GetTagByID(userID, tagID)
	if err != nil {
		handleServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse("Tag retrieved successfully", tag))
}

// @Summary Update tag
// @Description Update tag details
// @Tags Tags
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Tag ID"
// @Param request body dto.UpdateTagRequest true "Update Tag Request"
// @Success 200 {object} dto.Response
// @Failure 400 {object} dto.Response
// @Failure 401 {object} dto.Response
// @Failure 404 {object} dto.Response
// @Router /tags/{id} [put]
func (h *TagHandler) UpdateTag(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse("Invalid user ID in token"))
		return
	}

	tagIDStr := c.Param("id")
	tagID, err := uuid.Parse(tagIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse("Invalid tag ID format"))
		return
	}

	var req dto.UpdateTagRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse(err.Error()))
		return
	}

	tag, err := h.tagService.UpdateTag(userID, tagID, req)
	if err != nil {
		handleServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse("Tag updated successfully", tag))
}

// @Summary Delete tag
// @Description Delete a tag
// @Tags Tags
// @Produce json
// @Security BearerAuth
// @Param id path string true "Tag ID"
// @Success 200 {object} dto.Response
// @Failure 400 {object} dto.Response
// @Failure 401 {object} dto.Response
// @Failure 404 {object} dto.Response
// @Router /tags/{id} [delete]
func (h *TagHandler) DeleteTag(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse("Invalid user ID in token"))
		return
	}

	tagIDStr := c.Param("id")
	tagID, err := uuid.Parse(tagIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse("Invalid tag ID format"))
		return
	}

	if err := h.tagService.DeleteTag(userID, tagID); err != nil {
		handleServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse("Tag deleted successfully", nil))
}
