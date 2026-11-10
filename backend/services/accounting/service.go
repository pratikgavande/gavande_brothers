package accounting

import (
	"time"

	"billing-app/models"

	"gorm.io/gorm"
)

type Filters struct {
	Date   *time.Time
	Status *int8
}

type Service struct {
	db *gorm.DB
}

func NewService(db *gorm.DB) *Service {
	return &Service{db: db}
}

func (s *Service) List(page, limit int, filters Filters) ([]models.Accounting, int64, error) {
	query := s.db.Model(&models.Accounting{})
	if filters.Date != nil {
		query = query.Where("`date` = ?", filters.Date.Format("2006-01-02"))
	}
	if filters.Status != nil {
		query = query.Where("`status` = ?", *filters.Status)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	records := make([]models.Accounting, 0)
	err := query.Order("`date` DESC, `accId` DESC").
		Offset((page - 1) * limit).
		Limit(limit).
		Find(&records).Error
	return records, total, err
}

func (s *Service) Get(id int32) (models.Accounting, error) {
	var record models.Accounting
	err := s.db.First(&record, id).Error
	return record, err
}
