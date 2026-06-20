package handler

import (
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
// @Param request body service.CreateTagRequest true "Create Tag Request"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /api/v1/tags [post]
func (h *TagHandler) CreateTag(c *gin.Context) {
	userIdStr, _ := c.Get("userId")
	userID, err := uuid.Parse(userIdStr.(string))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID in token"})
		return
	}

	var req service.CreateTagRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tag, err := h.tagService.CreateTag(userID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Tag created successfully", "data": tag})
}

// @Summary Get all tags
// @Description Get all tags for the authenticated user
// @Tags Tags
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /api/v1/tags [get]
func (h *TagHandler) GetTags(c *gin.Context) {
	userIdStr, _ := c.Get("userId")
	userID, err := uuid.Parse(userIdStr.(string))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID in token"})
		return
	}

	tags, err := h.tagService.GetTags(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": tags})
}

// @Summary Get tag by ID
// @Description Get tag details by tag ID
// @Tags Tags
// @Produce json
// @Security BearerAuth
// @Param id path string true "Tag ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /api/v1/tags/{id} [get]
func (h *TagHandler) GetTagByID(c *gin.Context) {
	userIdStr, _ := c.Get("userId")
	userID, err := uuid.Parse(userIdStr.(string))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID in token"})
		return
	}

	tagIDStr := c.Param("id")
	tagID, err := uuid.Parse(tagIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tag ID format"})
		return
	}

	tag, err := h.tagService.GetTagByID(userID, tagID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": tag})
}

// @Summary Update tag
// @Description Update tag details
// @Tags Tags
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Tag ID"
// @Param request body service.UpdateTagRequest true "Update Tag Request"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /api/v1/tags/{id} [put]
func (h *TagHandler) UpdateTag(c *gin.Context) {
	userIdStr, _ := c.Get("userId")
	userID, err := uuid.Parse(userIdStr.(string))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID in token"})
		return
	}

	tagIDStr := c.Param("id")
	tagID, err := uuid.Parse(tagIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tag ID format"})
		return
	}

	var req service.UpdateTagRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tag, err := h.tagService.UpdateTag(userID, tagID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Tag updated successfully", "data": tag})
}

// @Summary Delete tag
// @Description Delete a tag
// @Tags Tags
// @Produce json
// @Security BearerAuth
// @Param id path string true "Tag ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /api/v1/tags/{id} [delete]
func (h *TagHandler) DeleteTag(c *gin.Context) {
	userIdStr, _ := c.Get("userId")
	userID, err := uuid.Parse(userIdStr.(string))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID in token"})
		return
	}

	tagIDStr := c.Param("id")
	tagID, err := uuid.Parse(tagIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tag ID format"})
		return
	}

	if err := h.tagService.DeleteTag(userID, tagID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Tag deleted successfully"})
}
