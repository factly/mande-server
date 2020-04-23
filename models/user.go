package models

import (
	"time"
)

// User model
type User struct {
	ID        uint      `gorm:"primary_key" json:"id"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`
	Email     string    `gorm:"column:email" json:"email" validate:"required"`
	FirstName string    `gorm:"column:first_name" json:"first_name" validate:"required"`
	LastName  string    `gorm:"column:last_name" json:"last_name" validate:"required"`
}
