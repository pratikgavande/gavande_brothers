package bill_items

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
	itemservice "billing-app/services/bill_items"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Controller struct {
	service *itemservice.Service
}

func NewController(db *gorm.DB) *Controller {
	return &Controller{service: itemservice.NewService(db)}
}

func (c *Controller) ListByBill(ctx *gin.Context) {
	billID, ok := response.ParseID(ctx, "id")
	if !ok {
		return
	}
	items, err := c.service.ListByBillID(billID)
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			log.Printf("list bill items failed: %v", err)
		}
		response.WriteDatabaseError(ctx, err, "could not retrieve bill items")
		return
	}
	response.Success(ctx, "Bill items retrieved successfully", items, nil)
}

func (c *Controller) Get(ctx *gin.Context) {
	id, ok := response.ParseID(ctx, "id")
	if !ok {
		return
	}
	item, err := c.service.Get(id)
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			log.Printf("get bill item failed: %v", err)
		}
		response.WriteDatabaseError(ctx, err, "could not retrieve bill item")
		return
	}
	response.Success(ctx, "Bill item retrieved successfully", item, nil)
}

type createBillItemRequest struct {
	Name   *string `json:"name"`
	Rate   *string `json:"rate"`
	Unit   *string `json:"unit"`
	Amount *int32  `json:"amount"`
	Date   *string `json:"date"`
	Status *int8   `json:"status"`
}

func (r createBillItemRequest) toModel() (models.BillItem, bool) {
	if r.Name == nil || r.Rate == nil || r.Unit == nil ||
		r.Amount == nil || r.Status == nil || strings.TrimSpace(*r.Name) == "" ||
		utf8.RuneCountInString(*r.Name) > 250 {
		return models.BillItem{}, false
	}
	return models.BillItem{
		Name:   strings.TrimSpace(*r.Name),
		Rate:   *r.Rate,
		Unit:   *r.Unit,
		Amount: *r.Amount,
		Date:   r.Date,
		Status: *r.Status,
	}, true
}

func (c *Controller) CreateMany(ctx *gin.Context) {
	billID, ok := response.ParseID(ctx, "id")
	if !ok {
		return
	}

	var request struct {
		Items []createBillItemRequest `json:"items"`
	}
	if err := ctx.ShouldBindJSON(&request); err != nil || len(request.Items) == 0 {
		response.Failure(ctx, http.StatusBadRequest, "items must be a non-empty array")
		return
	}

	items := make([]models.BillItem, len(request.Items))
	for i, itemRequest := range request.Items {
		item, valid := itemRequest.toModel()
		if !valid {
			response.Failure(ctx, http.StatusBadRequest, "each item requires name, rate, unit, amount, and status")
			return
		}
		items[i] = item
	}

	if err := c.service.CreateMany(billID, items); err != nil {
		if errors.Is(err, itemservice.ErrBillNotFound) {
			response.Failure(ctx, http.StatusNotFound, "bill not found")
			return
		}
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			response.Failure(ctx, http.StatusConflict, "a generated bill item ID already exists")
			return
		}
		log.Printf("create bill items failed: %v", err)
		response.Failure(ctx, http.StatusInternalServerError, "could not create bill items")
		return
	}
	response.Created(ctx, "Bill items created successfully", items)
}

func (c *Controller) Update(ctx *gin.Context) {
	id, ok := response.ParseID(ctx, "id")
	if !ok {
		return
	}

	var fields map[string]json.RawMessage
	if err := ctx.ShouldBindJSON(&fields); err != nil || len(fields) == 0 {
		response.Failure(ctx, http.StatusBadRequest, "provide at least one permitted bill item field to update")
		return
	}
	updates, err := billItemPatch(fields)
	if err != nil {
		response.Failure(ctx, http.StatusBadRequest, err.Error())
		return
	}

	item, err := c.service.Update(id, updates)
	if err != nil {
		if errors.Is(err, itemservice.ErrBillItemNotFound) {
			response.Failure(ctx, http.StatusNotFound, "bill item not found")
			return
		}
		log.Printf("update bill item failed: %v", err)
		response.Failure(ctx, http.StatusInternalServerError, "could not update bill item")
		return
	}
	response.Success(ctx, "Bill item updated successfully", item, nil)
}

func billItemPatch(fields map[string]json.RawMessage) (map[string]any, error) {
	columns := map[string]string{
		"name": "name", "rate": "rate", "unit": "unit",
		"amount": "amount", "date": "date", "status": "status",
	}
	updates := make(map[string]any, len(fields))
	for key, raw := range fields {
		column, ok := columns[key]
		if !ok {
			return nil, errors.New("unsupported bill item field: " + key)
		}
		if bytes.Equal(bytes.TrimSpace(raw), []byte("null")) {
			if key == "date" {
				updates[column] = nil
				continue
			}
			return nil, errors.New(key + " cannot be null")
		}
		switch key {
		case "name", "rate", "unit":
			var value string
			if err := json.Unmarshal(raw, &value); err != nil {
				return nil, errors.New(key + " must be a string")
			}
			if key == "name" {
				value = strings.TrimSpace(value)
				if value == "" || utf8.RuneCountInString(value) > 250 {
					return nil, errors.New("name must be non-empty and at most 250 characters")
				}
			}
			updates[column] = value
		case "amount":
			var value int32
			if err := json.Unmarshal(raw, &value); err != nil {
				return nil, errors.New("amount must be a 32-bit integer")
			}
			updates[column] = value
		case "status":
			var value int8
			if err := json.Unmarshal(raw, &value); err != nil {
				return nil, errors.New("status must be a signed TINYINT value")
			}
			updates[column] = value
		case "date":
			var value string
			if err := json.Unmarshal(raw, &value); err != nil {
				return nil, errors.New("date must be a string or null")
			}
			updates[column] = value
		}
	}
	return updates, nil
}
