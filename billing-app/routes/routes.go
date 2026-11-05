package routes

import (
	"billing-app/controllers"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Register(router *gin.Engine, db *gorm.DB) {
	healthController := controllers.NewHealthController(db)
	router.GET("/api/v1/health", healthController.Check)
}
