package roles

import (
	"errors"
	"log"
	"net/http"
	"strconv"

	roleservice "billing-app/services/roles"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Controller struct {
	service *roleservice.Service
}

func NewController(db *gorm.DB) *Controller {
	return &Controller{service: roleservice.NewService(db)}
}

func (c *Controller) List(ctx *gin.Context) {
	roles, err := c.service.List()
	if err != nil {
		log.Printf("load roles failed: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "could not load roles"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"roles": roles})
}

func (c *Controller) Create(ctx *gin.Context) {
	var request struct {
		Name        string `json:"name" binding:"required"`
		Code        string `json:"code" binding:"required"`
		Description string `json:"description"`
	}
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "name and code are required"})
		return
	}

	role, err := c.service.Create(request.Name, request.Code, request.Description)
	if errors.Is(err, roleservice.ErrInvalidRole) {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid role name, code, or description"})
		return
	}
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			ctx.JSON(http.StatusConflict, gin.H{"error": "role name or code already exists"})
			return
		}
		log.Printf("create role failed: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "could not create role"})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"role": role})
}

func (c *Controller) Update(ctx *gin.Context) {
	id, ok := roleID(ctx)
	if !ok {
		return
	}

	var request struct {
		Name        string `json:"name" binding:"required"`
		Description string `json:"description"`
	}
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "name is required"})
		return
	}

	role, err := c.service.Update(id, request.Name, request.Description)
	if err != nil {
		writeRoleError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"role": role})
}

func (c *Controller) Delete(ctx *gin.Context) {
	id, ok := roleID(ctx)
	if !ok {
		return
	}

	if err := c.service.Delete(id); err != nil {
		writeRoleError(ctx, err)
		return
	}
	ctx.Status(http.StatusNoContent)
}

func roleID(ctx *gin.Context) (uint64, bool) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil || id == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid role id"})
		return 0, false
	}
	return id, true
}

func writeRoleError(ctx *gin.Context, err error) {
	switch {
	case errors.Is(err, roleservice.ErrInvalidRole):
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid role name or description"})
	case errors.Is(err, roleservice.ErrRoleNotFound), errors.Is(err, gorm.ErrRecordNotFound):
		ctx.JSON(http.StatusNotFound, gin.H{"error": "role not found"})
	case errors.Is(err, roleservice.ErrRoleProtected):
		ctx.JSON(http.StatusConflict, gin.H{"error": "system roles cannot be changed or deleted"})
	case errors.Is(err, roleservice.ErrRoleInUse):
		ctx.JSON(http.StatusConflict, gin.H{"error": "role is assigned to users"})
	case errors.Is(err, gorm.ErrDuplicatedKey):
		ctx.JSON(http.StatusConflict, gin.H{"error": "role name already exists"})
	default:
		log.Printf("role operation failed: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "could not update role"})
	}
}
