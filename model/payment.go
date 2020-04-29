package model

import (
	"time"
)

// Payment model
type Payment struct {
	ID         uint      `gorm:"primary_key" json:"id"`
	CreatedAt  time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt  time.Time `gorm:"column:updated_at" json:"updated_at"`
	Amount     int       `gorm:"column:amount" json:"amount" validate:"required"`
	Gateway    string    `gorm:"column:gateway" json:"gateway" validate:"required"`
	CurrencyID uint      `gorm:"column:currency_id" json:"currency_id" validate:"required"`
	Currency   Currency  `gorm:"foreignkey:currency_id;association_foreignkey:id"`
	Status     string    `gorm:"column:status" json:"status" validate:"required"`
}
