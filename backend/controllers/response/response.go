package response

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

const (
	defaultPage  = 1
	defaultLimit = 10
	maxLimit     = 100
)

type Pagination struct {
	Page  int   `json:"page"`
	Limit int   `json:"limit"`
	Total int64 `json:"total"`
}

func ParsePagination(c *gin.Context) (int, int, bool) {
	page, err := parseQueryInt(c, "page", defaultPage)
	if err != nil || page < 1 {
		Failure(c, http.StatusBadRequest, "page must be a positive integer")
		return 0, 0, false
	}

	limit, err := parseQueryInt(c, "limit", defaultLimit)
	if err != nil || limit < 1 || limit > maxLimit {
		Failure(c, http.StatusBadRequest, "limit must be between 1 and 100")
		return 0, 0, false
	}
	if page > int(^uint(0)>>1)/limit {
		Failure(c, http.StatusBadRequest, "page is too large")
		return 0, 0, false
	}
	return page, limit, true
}

func Success(c *gin.Context, message string, data any, pagination *Pagination) {
	writeSuccess(c, http.StatusOK, message, data, pagination)
}

func Created(c *gin.Context, message string, data any) {
	writeSuccess(c, http.StatusCreated, message, data, nil)
}

func writeSuccess(c *gin.Context, status int, message string, data any, pagination *Pagination) {
	body := gin.H{
		"success": true,
		"message": message,
		"data":    data,
	}
	if pagination != nil {
		body["pagination"] = pagination
	}
	c.JSON(status, body)
}

func Failure(c *gin.Context, status int, message string) {
	c.JSON(status, gin.H{"success": false, "message": message})
}

func ParseID(c *gin.Context, param string) (int32, bool) {
	value, err := strconv.ParseInt(c.Param(param), 10, 32)
	if err != nil || value < 1 {
		Failure(c, http.StatusBadRequest, "id must be a positive 32-bit integer")
		return 0, false
	}
	return int32(value), true
}

func WriteDatabaseError(c *gin.Context, err error, message string) {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		Failure(c, http.StatusNotFound, "record not found")
		return
	}
	Failure(c, http.StatusInternalServerError, message)
}

func parseQueryInt(c *gin.Context, key string, fallback int) (int, error) {
	value := c.Query(key)
	if value == "" {
		return fallback, nil
	}
	return strconv.Atoi(value)
}
