package bill_items

import (
	"errors"

	"billing-app/models"

	"gorm.io/gorm"
)

var ErrBillNotFound = errors.New("bill not found")
var ErrBillItemNotFound = errors.New("bill item not found")

type Service struct {
	db *gorm.DB
}

func NewService(db *gorm.DB) *Service {
	return &Service{db: db}
}

func (s *Service) ListByBillID(billID int32) ([]models.BillItem, error) {
	var bill models.Bill
	if err := s.db.Select("billId").First(&bill, billID).Error; err != nil {
		return nil, err
	}

	items := make([]models.BillItem, 0)
	err := s.db.Where("`billId` = ?", billID).
		Order("`itemId` ASC").
		Find(&items).Error
	return items, err
}

func (s *Service) Get(id int32) (models.BillItem, error) {
	var item models.BillItem
	err := s.db.First(&item, id).Error
	return item, err
}

func (s *Service) CreateMany(billID int32, items []models.BillItem) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		var bill models.Bill
		if err := tx.Select("billId").First(&bill, billID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrBillNotFound
			}
			return err
		}

		for i := range items {
			items[i].BillID = billID
		}
		return tx.Create(&items).Error
	})
}

func (s *Service) Update(id int32, fields map[string]any) (models.BillItem, error) {
	var item models.BillItem
	err := s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.First(&item, id).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrBillItemNotFound
			}
			return err
		}
		if err := tx.Model(&item).Updates(fields).Error; err != nil {
			return err
		}
		return tx.First(&item, id).Error
	})
	return item, err
}
