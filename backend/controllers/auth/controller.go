package auth

import (
	"errors"
	"log"
	"net/http"

	"billing-app/middleware"
	authservice "billing-app/services/auth"

	"github.com/gin-gonic/gin"
)

type Controller struct {
	service *authservice.Service
}

func NewController(service *authservice.Service) *Controller {
	return &Controller{service: service}
}

func (c *Controller) Login(ctx *gin.Context) {
	var request struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "valid email and password are required"})
		return
	}

	result, err := c.service.Login(request.Email, request.Password)
	if errors.Is(err, authservice.ErrInvalidCredentials) {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid email or password"})
		return
	}
	if err != nil {
		log.Printf("login failed: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "could not complete login"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"access_token": result.AccessToken,
		"token_type":   "Bearer",
		"expires_at":   result.ExpiresAt,
		"user":         result.User,
	})
}

func (c *Controller) Me(ctx *gin.Context) {
	user, ok := middleware.CurrentUser(ctx)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "authentication required"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"user": user})
}
