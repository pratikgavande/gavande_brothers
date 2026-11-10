package models

import "time"

type User struct {
	ID           uint64    `json:"id" gorm:"primaryKey;autoIncrement"`
	FullName     string    `json:"full_name" gorm:"size:150;not null"`
	Email        string    `json:"email" gorm:"size:191;not null;uniqueIndex"`
	PasswordHash string    `json:"-" gorm:"size:255;not null"`
	RoleID       uint64    `json:"role_id" gorm:"not null;index"`
	Role         Role      `json:"role" gorm:"foreignKey:RoleID"`
	IsActive     bool      `json:"is_active" gorm:"not null;default:true"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func (User) TableName() string {
	return "users"
}
