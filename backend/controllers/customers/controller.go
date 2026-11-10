package customers

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"
	"unicode/utf8"

	"billing-app/controllers/response"
	"billing-app/models"
	customerservice "billing-app/services/customers"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Controller struct {
	service *customerservice.Service
}

func NewController(db *gorm.DB) *Controller {
	return &Controller{service: customerservice.NewService(db)}
}

func (c *Controller) List(ctx *gin.Context) {
	page, limit, ok := response.ParsePagination(ctx)
	if !ok {
		return
	}
	search := strings.TrimSpace(ctx.Query("search"))
	if len(search) > 250 {
		response.Failure(ctx, http.StatusBadRequest, "search must not exceed 250 characters")
		return
	}

	customers, total, err := c.service.List(page, limit, search)
	if err != nil {
		log.Printf("list customers failed: %v", err)
		response.Failure(ctx, http.StatusInternalServerError, "could not retrieve customers")
		return
	}
	response.Success(ctx, "Customers retrieved successfully", customers, &response.Pagination{
		Page: page, Limit: limit, Total: total,
	})
}

func (c *Controller) Get(ctx *gin.Context) {
	id, ok := response.ParseID(ctx, "id")
	if !ok {
		return
	}
	customer, err := c.service.Get(id)
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			log.Printf("get customer failed: %v", err)
		}
		response.WriteDatabaseError(ctx, err, "could not retrieve customer")
		return
	}
	response.Success(ctx, "Customer retrieved successfully", customer, nil)
}

type createCustomerRequest struct {
	Name    *string `json:"name"`
	Mobile  *string `json:"mobile"`
	GSTNo   *string `json:"gst_no"`
	Email   *string `json:"email"`
	Address *string `json:"address"`
	Status  *int8   `json:"status"`
}

func (r createCustomerRequest) toModel() (models.Customer, bool) {
	if r.Name == nil || r.Mobile == nil || r.GSTNo == nil ||
		r.Email == nil || r.Address == nil || r.Status == nil ||
		strings.TrimSpace(*r.Name) == "" || strings.TrimSpace(*r.Mobile) == "" ||
		utf8.RuneCountInString(*r.Name) > 250 ||
		utf8.RuneCountInString(*r.Mobile) > 250 ||
		utf8.RuneCountInString(*r.GSTNo) > 250 ||
		utf8.RuneCountInString(*r.Email) > 250 {
		return models.Customer{}, false
	}
	return models.Customer{
		Name:    *r.Name,
		Mobile:  *r.Mobile,
		GSTNo:   *r.GSTNo,
		Email:   *r.Email,
		Address: *r.Address,
		Status:  *r.Status,
	}, true
}

func (c *Controller) Create(ctx *gin.Context) {
	var request createCustomerRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		response.Failure(ctx, http.StatusBadRequest, "invalid customer request")
		return
	}
	customer, valid := request.toModel()
	if !valid {
		response.Failure(ctx, http.StatusBadRequest, "all customer fields are required, name and mobile must be non-empty, and varchar fields must not exceed 250 characters")
		return
	}

	if err := c.service.Create(&customer); err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			response.Failure(ctx, http.StatusConflict, "customer conflicts with an existing unique value")
			return
		}
		log.Printf("create customer failed: %v", err)
		response.Failure(ctx, http.StatusInternalServerError, "could not create customer")
		return
	}
	response.Created(ctx, "Customer created successfully", customer)
}

func (c *Controller) Update(ctx *gin.Context) {
	id, ok := response.ParseID(ctx, "id")
	if !ok {
		return
	}

	var fields map[string]json.RawMessage
	if err := ctx.ShouldBindJSON(&fields); err != nil || len(fields) == 0 {
		response.Failure(ctx, http.StatusBadRequest, "provide at least one customer field to update")
		return
	}
	updates, err := customerPatch(fields)
	if err != nil {
		response.Failure(ctx, http.StatusBadRequest, err.Error())
		return
	}

	customer, err := c.service.Update(id, updates)
	if err != nil {
		if errors.Is(err, customerservice.ErrCustomerNotFound) {
			response.Failure(ctx, http.StatusNotFound, "customer not found")
			return
		}
		log.Printf("update customer failed: %v", err)
		response.Failure(ctx, http.StatusInternalServerError, "could not update customer")
		return
	}
	response.Success(ctx, "Customer updated successfully", customer, nil)
}

func customerPatch(fields map[string]json.RawMessage) (map[string]any, error) {
	columns := map[string]string{
		"name": "name", "mobile": "mobile", "gst_no": "gst_no",
		"email": "email", "address": "address",
	}
	updates := make(map[string]any, len(fields))
	for key, raw := range fields {
		column, ok := columns[key]
		if !ok {
			return nil, errors.New("unsupported customer field: " + key)
		}
		if bytes.Equal(bytes.TrimSpace(raw), []byte("null")) {
			return nil, errors.New(key + " cannot be null")
		}
		var value string
		if err := json.Unmarshal(raw, &value); err != nil {
			return nil, errors.New(key + " must be a string")
		}
		if (key == "name" || key == "mobile") && strings.TrimSpace(value) == "" {
			return nil, errors.New(key + " must not be blank")
		}
		if key != "address" && utf8.RuneCountInString(value) > 250 {
			return nil, errors.New(key + " must not exceed 250 characters")
		}
		updates[column] = value
	}
	return updates, nil
}
