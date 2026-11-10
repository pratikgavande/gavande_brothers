package bills

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"net/mail"
	"strconv"
	"time"

	"billing-app/config"
	"billing-app/controllers/response"
	"billing-app/models"
	"billing-app/services/bill_documents"
	"billing-app/services/bill_email"
	billservice "billing-app/services/bills"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Controller struct {
	service   *billservice.Service
	documents *bill_documents.Service
	billEmail *bill_email.Service
}

func NewController(db *gorm.DB, smtpConfig config.SMTPConfig) *Controller {
	return &Controller{
		service:   billservice.NewService(db),
		documents: bill_documents.NewService(db),
		billEmail: bill_email.NewService(smtpConfig),
	}
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

	bills, total, err := c.service.List(page, limit, filters)
	if err != nil {
		log.Printf("list bills failed: %v", err)
		response.Failure(ctx, http.StatusInternalServerError, "could not retrieve bills")
		return
	}
	response.Success(ctx, "Bills retrieved successfully", bills, &response.Pagination{
		Page: page, Limit: limit, Total: total,
	})
}

func (c *Controller) Get(ctx *gin.Context) {
	id, ok := response.ParseID(ctx, "id")
	if !ok {
		return
	}
	bill, err := c.service.Get(id)
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			log.Printf("get bill failed: %v", err)
		}
		response.WriteDatabaseError(ctx, err, "could not retrieve bill")
		return
	}
	response.Success(ctx, "Bill retrieved successfully", bill, nil)
}

type createBillRequest struct {
	BillNo         *int64     `json:"bill_no"`
	CustomerID     *int32     `json:"customer_id"`
	Date           *time.Time `json:"date"`
	GSTStatus      *int8      `json:"gst_status"`
	GSTAmount      *int32     `json:"gst_amount"`
	Amount         *int32     `json:"amount"`
	AmountReceived *int32     `json:"amount_received"`
	AmountPending  *int32     `json:"amount_pending"`
	EngineersName  *string    `json:"engineers_name"`
	SiteName       *string    `json:"site_name"`
	Location       *string    `json:"location"`
	ItemCount      *int32     `json:"item_count"`
	Status         *int8      `json:"status"`
}

func (r createBillRequest) toModel() (models.Bill, bool) {
	if r.BillNo == nil || r.CustomerID == nil || r.Date == nil ||
		r.GSTStatus == nil || r.GSTAmount == nil || r.Amount == nil ||
		r.AmountReceived == nil || r.AmountPending == nil ||
		r.EngineersName == nil || r.SiteName == nil || r.Location == nil ||
		r.ItemCount == nil || r.Status == nil || *r.CustomerID < 1 {
		return models.Bill{}, false
	}

	return models.Bill{
		BillNo:         *r.BillNo,
		CustomerID:     *r.CustomerID,
		Date:           *r.Date,
		GSTStatus:      *r.GSTStatus,
		GSTAmount:      *r.GSTAmount,
		Amount:         *r.Amount,
		AmountReceived: *r.AmountReceived,
		AmountPending:  *r.AmountPending,
		EngineersName:  *r.EngineersName,
		SiteName:       *r.SiteName,
		Location:       *r.Location,
		ItemCount:      *r.ItemCount,
		Status:         *r.Status,
	}, true
}

func (c *Controller) Create(ctx *gin.Context) {
	var request createBillRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		response.Failure(ctx, http.StatusBadRequest, "invalid bill request; date must be RFC3339")
		return
	}
	bill, valid := request.toModel()
	if !valid {
		response.Failure(ctx, http.StatusBadRequest, "all bill fields are required and customer_id must be positive")
		return
	}

	if err := c.service.Create(&bill); err != nil {
		if errors.Is(err, billservice.ErrCustomerNotFound) {
			response.Failure(ctx, http.StatusNotFound, "customer not found")
			return
		}
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			response.Failure(ctx, http.StatusConflict, "bill conflicts with an existing unique value")
			return
		}
		log.Printf("create bill failed: %v", err)
		response.Failure(ctx, http.StatusInternalServerError, "could not create bill")
		return
	}
	response.Created(ctx, "Bill created successfully", bill)
}

func (c *Controller) Update(ctx *gin.Context) {
	id, ok := response.ParseID(ctx, "id")
	if !ok {
		return
	}

	var fields map[string]json.RawMessage
	if err := ctx.ShouldBindJSON(&fields); err != nil || len(fields) == 0 {
		response.Failure(ctx, http.StatusBadRequest, "provide at least one permitted bill field to update")
		return
	}
	updates, err := billPatch(fields)
	if err != nil {
		response.Failure(ctx, http.StatusBadRequest, err.Error())
		return
	}

	bill, err := c.service.Update(id, updates)
	if err != nil {
		if errors.Is(err, billservice.ErrBillNotFound) {
			response.Failure(ctx, http.StatusNotFound, "bill not found")
			return
		}
		log.Printf("update bill failed: %v", err)
		response.Failure(ctx, http.StatusInternalServerError, "could not update bill")
		return
	}
	response.Success(ctx, "Bill updated successfully", bill, nil)
}

func (c *Controller) DownloadPDF(ctx *gin.Context) {
	id, ok := response.ParseID(ctx, "id")
	if !ok {
		return
	}
	document, err := c.documents.Load(id)
	if err != nil {
		c.writeDocumentError(ctx, err, "could not generate bill PDF")
		return
	}
	pdf, err := bill_documents.RenderPDF(document)
	if err != nil {
		log.Printf("render bill PDF failed: %v", err)
		response.Failure(ctx, http.StatusInternalServerError, "could not generate bill PDF")
		return
	}
	filename := "gavande-bros-bill-" + strconv.FormatInt(document.Bill.BillNo, 10) + ".pdf"
	ctx.Header("Content-Type", "application/pdf")
	ctx.Header("Content-Disposition", `attachment; filename="`+filename+`"`)
	ctx.Data(http.StatusOK, "application/pdf", pdf)
}

func (c *Controller) Send(ctx *gin.Context) {
	id, ok := response.ParseID(ctx, "id")
	if !ok {
		return
	}
	document, err := c.documents.Load(id)
	if err != nil {
		c.writeDocumentError(ctx, err, "could not send bill")
		return
	}
	if _, err := mail.ParseAddress(document.Customer.Email); err != nil {
		response.Failure(ctx, http.StatusUnprocessableEntity, "customer does not have a valid email address")
		return
	}
	pdf, err := bill_documents.RenderPDF(document)
	if err != nil {
		log.Printf("render bill PDF for email failed: %v", err)
		response.Failure(ctx, http.StatusInternalServerError, "could not generate bill PDF")
		return
	}
	sendCtx, cancel := context.WithTimeout(ctx.Request.Context(), 30*time.Second)
	defer cancel()
	if err := c.billEmail.SendBill(sendCtx, document.Customer.Email, document.Bill.BillNo, pdf); err != nil {
		if errors.Is(err, bill_email.ErrNotConfigured) {
			response.Failure(ctx, http.StatusServiceUnavailable, "email service is not configured")
			return
		}
		if errors.Is(err, bill_email.ErrInvalidEmail) {
			response.Failure(ctx, http.StatusUnprocessableEntity, "customer does not have a valid email address")
			return
		}
		log.Printf("send bill email failed: %T", err)
		response.Failure(ctx, http.StatusBadGateway, "email provider did not confirm delivery")
		return
	}
	response.Success(ctx, "Bill email accepted by SMTP server", gin.H{
		"bill_id": document.Bill.BillID,
	}, nil)
}

func (c *Controller) writeDocumentError(ctx *gin.Context, err error, message string) {
	switch {
	case errors.Is(err, bill_documents.ErrBillNotFound):
		response.Failure(ctx, http.StatusNotFound, "bill not found")
	case errors.Is(err, bill_documents.ErrCustomerNotFound):
		response.Failure(ctx, http.StatusNotFound, "customer associated with bill not found")
	default:
		log.Printf("%s: %v", message, err)
		response.Failure(ctx, http.StatusInternalServerError, message)
	}
}

func billPatch(fields map[string]json.RawMessage) (map[string]any, error) {
	columns := map[string]string{
		"date": "date", "gst_status": "gstStatus", "gst_amount": "gstAmount",
		"engineers_name": "engineersName", "site_name": "siteName",
		"location": "location", "status": "status",
	}
	updates := make(map[string]any, len(fields))
	for key, raw := range fields {
		column, ok := columns[key]
		if !ok {
			return nil, errors.New("unsupported bill field: " + key)
		}
		if bytes.Equal(bytes.TrimSpace(raw), []byte("null")) {
			return nil, errors.New(key + " cannot be null")
		}
		switch key {
		case "date":
			var value time.Time
			if err := json.Unmarshal(raw, &value); err != nil {
				return nil, errors.New("date must be an RFC3339 date-time")
			}
			updates[column] = value
		case "gst_status", "status":
			var value int8
			if err := json.Unmarshal(raw, &value); err != nil {
				return nil, errors.New(key + " must be a signed TINYINT value")
			}
			updates[column] = value
		case "gst_amount":
			var value int32
			if err := json.Unmarshal(raw, &value); err != nil {
				return nil, errors.New("gst_amount must be a 32-bit integer")
			}
			updates[column] = value
		default:
			var value string
			if err := json.Unmarshal(raw, &value); err != nil {
				return nil, errors.New(key + " must be a string")
			}
			updates[column] = value
		}
	}
	return updates, nil
}

func parseFilters(ctx *gin.Context) (billservice.Filters, bool) {
	var filters billservice.Filters
	if raw := ctx.Query("customerId"); raw != "" {
		value, err := strconv.ParseInt(raw, 10, 32)
		if err != nil || value < 1 {
			response.Failure(ctx, http.StatusBadRequest, "customerId must be a positive 32-bit integer")
			return filters, false
		}
		customerID := int32(value)
		filters.CustomerID = &customerID
	}
	if raw := ctx.Query("billNo"); raw != "" {
		value, err := strconv.ParseInt(raw, 10, 64)
		if err != nil {
			response.Failure(ctx, http.StatusBadRequest, "billNo must be a 64-bit integer")
			return filters, false
		}
		filters.BillNo = &value
	}
	if raw := ctx.Query("date"); raw != "" {
		date, err := time.Parse("2006-01-02", raw)
		if err != nil {
			response.Failure(ctx, http.StatusBadRequest, "date must use YYYY-MM-DD format")
			return filters, false
		}
		filters.Date = &date
	}
	return filters, true
}
