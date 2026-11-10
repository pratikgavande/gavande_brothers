package bill_documents

import (
	"errors"

	"billing-app/models"

	"gorm.io/gorm"
)

var (
	ErrBillNotFound     = errors.New("bill not found")
	ErrCustomerNotFound = errors.New("customer not found")
)

type Document struct {
	Bill     models.Bill
	Customer models.Customer
	Items    []models.BillItem
}

type Service struct {
	db *gorm.DB
}

func NewService(db *gorm.DB) *Service {
	return &Service{db: db}
}

func (s *Service) Load(billID int32) (Document, error) {
	var document Document
	if err := s.db.First(&document.Bill, billID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return Document{}, ErrBillNotFound
		}
		return Document{}, err
	}
	if err := s.db.First(&document.Customer, document.Bill.CustomerID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return Document{}, ErrCustomerNotFound
		}
		return Document{}, err
	}

	document.Items = make([]models.BillItem, 0)
	if err := s.db.Where("`billId` = ?", document.Bill.BillID).
		Order("`itemId` ASC").
		Find(&document.Items).Error; err != nil {
		return Document{}, err
	}
	return document, nil
}
