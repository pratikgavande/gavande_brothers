package models

import "time"

type Role struct {
	ID          uint64    `json:"id" gorm:"primaryKey;autoIncrement"`
	Name        string    `json:"name" gorm:"size:100;not null;uniqueIndex"`
	Code        string    `json:"code" gorm:"size:100;not null;uniqueIndex"`
	Description string    `json:"description" gorm:"size:255;not null;default:''"`
	IsSystem    bool      `json:"is_system" gorm:"not null;default:false"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (Role) TableName() string {
	return "roles"
}
