package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type HealthController struct {
	db *gorm.DB
}

func NewHealthController(db *gorm.DB) *HealthController {
	return &HealthController{db: db}
}

func (h *HealthController) Check(c *gin.Context) {
	sqlDB, err := h.db.DB()
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"status": "unavailable"})
		return
	}

	if err := sqlDB.PingContext(c.Request.Context()); err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"status": "unavailable"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
