package customers

import (
	"errors"
	"strings"

	"billing-app/models"

	"gorm.io/gorm"
)

var ErrCustomerNotFound = errors.New("customer not found")

type Service struct {
	db *gorm.DB
}

func NewService(db *gorm.DB) *Service {
	return &Service{db: db}
}

func (s *Service) List(page, limit int, search string) ([]models.Customer, int64, error) {
	query := s.db.Model(&models.Customer{})
	if search != "" {
		pattern := "%" + strings.TrimSpace(search) + "%"
		query = query.Where("`name` LIKE ? OR `mobile` LIKE ?", pattern, pattern)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	customers := make([]models.Customer, 0)
	err := query.Order("`customerId` ASC").
		Offset((page - 1) * limit).
		Limit(limit).
		Find(&customers).Error
	return customers, total, err
}

func (s *Service) Get(id int32) (models.Customer, error) {
	var customer models.Customer
	err := s.db.First(&customer, id).Error
	return customer, err
}

func (s *Service) Create(customer *models.Customer) error {
	return s.db.Create(customer).Error
}

func (s *Service) Update(id int32, fields map[string]any) (models.Customer, error) {
	var customer models.Customer
	err := s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.First(&customer, id).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrCustomerNotFound
			}
			return err
		}
		if err := tx.Model(&customer).Updates(fields).Error; err != nil {
			return err
		}
		return tx.First(&customer, id).Error
	})
	return customer, err
}
