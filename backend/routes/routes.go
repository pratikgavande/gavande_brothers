package routes

import (
	"billing-app/config"
	"billing-app/controllers"
	accountingcontroller "billing-app/controllers/accounting"
	authcontroller "billing-app/controllers/auth"
	itemcontroller "billing-app/controllers/bill_items"
	billcontroller "billing-app/controllers/bills"
	customercontroller "billing-app/controllers/customers"
	rolecontroller "billing-app/controllers/roles"
	"billing-app/middleware"
	"billing-app/services/auth"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Register(router *gin.Engine, db *gorm.DB, authService *auth.Service, jwtSecret []byte, smtpConfig config.SMTPConfig) {
	healthController := controllers.NewHealthController(db)
	router.GET("/api/v1/health", healthController.Check)

	authController := authcontroller.NewController(authService)
	roleController := rolecontroller.NewController(db)
	customerController := customercontroller.NewController(db)
	billController := billcontroller.NewController(db, smtpConfig)
	billItemController := itemcontroller.NewController(db)
	accountingController := accountingcontroller.NewController(db)

	api := router.Group("/api/v1")
	api.POST("/auth/login", authController.Login)

	protected := api.Group("")
	protected.Use(middleware.Authenticate(db, jwtSecret))
	protected.GET("/auth/me", authController.Me)
	protected.GET("/customers", customerController.List)
	protected.GET("/customers/:id", customerController.Get)
	protected.POST("/customers", customerController.Create)
	protected.GET("/bills", billController.List)
	protected.GET("/bills/:id", billController.Get)
	protected.POST("/bills", billController.Create)
	protected.PATCH("/bills/:id", billController.Update)
	protected.GET("/bills/:id/pdf", billController.DownloadPDF)
	protected.POST("/bills/:id/send", billController.Send)
	protected.GET("/bills/:id/items", billItemController.ListByBill)
	protected.POST("/bills/:id/items", billItemController.CreateMany)
	protected.GET("/bill-items/:id", billItemController.Get)
	protected.PATCH("/bill-items/:id", billItemController.Update)
	protected.GET("/accounting", accountingController.List)
	protected.GET("/accounting/:id", accountingController.Get)

	superAdmin := protected.Group("")
	superAdmin.Use(middleware.RequireSuperAdmin())
	superAdmin.GET("/roles", roleController.List)
	superAdmin.POST("/roles", roleController.Create)
	superAdmin.PATCH("/roles/:id", roleController.Update)
	superAdmin.DELETE("/roles/:id", roleController.Delete)
}
