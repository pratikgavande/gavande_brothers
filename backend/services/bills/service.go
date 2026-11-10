package bills

import (
	"errors"
	"time"

	"billing-app/models"

	"gorm.io/gorm"
)

var ErrCustomerNotFound = errors.New("customer not found")
var ErrBillNotFound = errors.New("bill not found")

type Filters struct {
	CustomerID *int32
	Date       *time.Time
	BillNo     *int64
}

type Service struct {
	db *gorm.DB
}

func NewService(db *gorm.DB) *Service {
	return &Service{db: db}
}

func (s *Service) List(page, limit int, filters Filters) ([]models.Bill, int64, error) {
	query := s.db.Model(&models.Bill{})
	if filters.CustomerID != nil {
		query = query.Where("`customerId` = ?", *filters.CustomerID)
	}
	if filters.Date != nil {
		query = query.Where("DATE(`date`) = ?", filters.Date.Format("2006-01-02"))
	}
	if filters.BillNo != nil {
		query = query.Where("`billNo` = ?", *filters.BillNo)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	bills := make([]models.Bill, 0)
	err := query.Order("`billId` DESC").
		Offset((page - 1) * limit).
		Limit(limit).
		Find(&bills).Error
	return bills, total, err
}

func (s *Service) Get(id int32) (models.Bill, error) {
	var bill models.Bill
	err := s.db.First(&bill, id).Error
	return bill, err
}

func (s *Service) Create(bill *models.Bill) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		var customer models.Customer
		if err := tx.Select("customerId").First(&customer, bill.CustomerID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrCustomerNotFound
			}
			return err
		}
		return tx.Create(bill).Error
	})
}

func (s *Service) Update(id int32, fields map[string]any) (models.Bill, error) {
	var bill models.Bill
	err := s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.First(&bill, id).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrBillNotFound
			}
			return err
		}
		if err := tx.Model(&bill).Updates(fields).Error; err != nil {
			return err
		}
		return tx.First(&bill, id).Error
	})
	return bill, err
}
