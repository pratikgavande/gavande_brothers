package accounting

import (
	"errors"
	"log"
	"net/http"
	"strconv"
	"time"

	"billing-app/controllers/response"
	accountingservice "billing-app/services/accounting"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Controller struct {
	service *accountingservice.Service
}

func NewController(db *gorm.DB) *Controller {
	return &Controller{service: accountingservice.NewService(db)}
}

func (c *Controller) List(ctx *gin.Context) {
	page, limit, ok := response.ParsePagination(ctx)
	if !ok {
		return
	}
	filters, ok := parseFilters(ctx)
	if !ok {
		return
	}

	records, total, err := c.service.List(page, limit, filters)
	if err != nil {
		log.Printf("list accounting records failed: %v", err)
		response.Failure(ctx, http.StatusInternalServerError, "could not retrieve accounting records")
		return
	}
	response.Success(ctx, "Accounting records retrieved successfully", records, &response.Pagination{
		Page: page, Limit: limit, Total: total,
	})
}

func (c *Controller) Get(ctx *gin.Context) {
	id, ok := response.ParseID(ctx, "id")
	if !ok {
		return
	}
	record, err := c.service.Get(id)
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			log.Printf("get accounting record failed: %v", err)
		}
		response.WriteDatabaseError(ctx, err, "could not retrieve accounting record")
		return
	}
	response.Success(ctx, "Accounting record retrieved successfully", record, nil)
}

func parseFilters(ctx *gin.Context) (accountingservice.Filters, bool) {
	var filters accountingservice.Filters
	if raw := ctx.Query("date"); raw != "" {
		date, err := time.Parse("2006-01-02", raw)
		if err != nil {
			response.Failure(ctx, http.StatusBadRequest, "date must use YYYY-MM-DD format")
			return filters, false
		}
		filters.Date = &date
	}
	if raw := ctx.Query("status"); raw != "" {
		value, err := strconv.ParseInt(raw, 10, 8)
		if err != nil {
			response.Failure(ctx, http.StatusBadRequest, "status must be an integer from -128 to 127")
			return filters, false
		}
		status := int8(value)
		filters.Status = &status
	}
	return filters, true
}
