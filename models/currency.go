package models

import (
	"time"
)

// Currency model
type Currency struct {
	ID        uint      `gorm:"primary_key" json:"id"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`
	IsoCode   string    `gorm:"column:iso_code" json:"iso_code" validate:"required"`
	Name      string    `gorm:"column:name" json:"name" validate:"required"`
}
