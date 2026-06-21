package dto

import (
	"strconv"
	"github.com/gin-gonic/gin"
)

type Response struct {
	Msg     string `json:"msg"`
	Success bool   `json:"success"`
	Data    any    `json:"data"`
	Error   string `json:"error,omitempty"`
	Meta    any    `json:"meta,omitempty"`
}

type PaginationMeta struct {
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	TotalItems int64 `json:"total_items"`
	TotalPages int64 `json:"total_pages"`
}

func SuccessResponse(msg string, data any) Response {
	return Response{
		Msg:     msg,
		Success: true,
		Data:    data,
	}
}

func PaginatedResponse(msg string, data any, page, limit int, totalItems int64) Response {
	limit64 := int64(limit)
	if limit64 <= 0 {
		limit64 = 10
	}
	totalPages := totalItems / limit64
	if totalItems%limit64 != 0 {
		totalPages++
	}
	if totalPages == 0 {
		totalPages = 1
	}

	return Response{
		Msg:     msg,
		Success: true,
		Data:    data,
		Meta: PaginationMeta{
			Page:       page,
			Limit:      limit,
			TotalItems: totalItems,
			TotalPages: totalPages,
		},
	}
}

func ErrorResponse(err string) Response {
	return Response{
		Msg:     "Failed",
		Success: false,
		Data:    nil,
		Error:   err,
	}
}

func GetPaginationParams(c *gin.Context) (int, int) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	if page < 1 {
		page = 1
	}
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if limit < 1 {
		limit = 10
	}
	return page, limit
}
