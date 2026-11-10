package roles

import (
	"errors"
	"regexp"
	"strings"

	"billing-app/models"

	"gorm.io/gorm"
)

var (
	ErrInvalidRole   = errors.New("invalid role")
	ErrRoleNotFound  = errors.New("role not found")
	ErrRoleProtected = errors.New("system role cannot be changed or deleted")
	ErrRoleInUse     = errors.New("role is assigned to users")
)

var roleCodePattern = regexp.MustCompile(`^[a-z][a-z0-9_]{0,49}$`)

type Service struct {
	db *gorm.DB
}

func NewService(db *gorm.DB) *Service {
	return &Service{db: db}
}

func (s *Service) List() ([]models.Role, error) {
	var roles []models.Role
	err := s.db.Order("name ASC").Find(&roles).Error
	return roles, err
}

func (s *Service) Create(name, code, description string) (models.Role, error) {
	name = strings.TrimSpace(name)
	code = strings.ToLower(strings.TrimSpace(code))
	description = strings.TrimSpace(description)
	if name == "" || len(name) > 100 || !roleCodePattern.MatchString(code) || len(description) > 255 {
		return models.Role{}, ErrInvalidRole
	}

	role := models.Role{Name: name, Code: code, Description: description}
	if err := s.db.Create(&role).Error; err != nil {
		return models.Role{}, err
	}
	return role, nil
}

func (s *Service) Update(id uint64, name, description string) (models.Role, error) {
	name = strings.TrimSpace(name)
	description = strings.TrimSpace(description)
	if name == "" || len(name) > 100 || len(description) > 255 {
		return models.Role{}, ErrInvalidRole
	}

	var role models.Role
	if err := s.db.First(&role, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.Role{}, ErrRoleNotFound
		}
		return models.Role{}, err
	}
	if role.IsSystem {
		return models.Role{}, ErrRoleProtected
	}

	role.Name = name
	role.Description = description
	if err := s.db.Save(&role).Error; err != nil {
		return models.Role{}, err
	}
	return role, nil
}

func (s *Service) Delete(id uint64) error {
	var role models.Role
	if err := s.db.First(&role, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrRoleNotFound
		}
		return err
	}
	if role.IsSystem {
		return ErrRoleProtected
	}

	var userCount int64
	if err := s.db.Model(&models.User{}).Where("role_id = ?", id).Count(&userCount).Error; err != nil {
		return err
	}
	if userCount > 0 {
		return ErrRoleInUse
	}

	if err := s.db.Delete(&role).Error; err != nil {
		if errors.Is(err, gorm.ErrForeignKeyViolated) {
			return ErrRoleInUse
		}
		return err
	}
	return nil
}
