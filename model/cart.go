package model

import (
	"time"
)

// Cart model
type Cart struct {
	ID        uint      `gorm:"primary_key" json:"id"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`
	Status    string    `gorm:"column:status" json:"status" validate:"required"`
	UserID    uint      `gorm:"column:user_id" json:"user_id" validate:"required"`
}

// CartItem model
type CartItem struct {
	ID        uint      `gorm:"primary_key" json:"id"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`
	IsDeleted bool      `gorm:"column:is_deleted" json:"is_deleted"`
	CartID    uint      `gorm:"column:cart_id" json:"cart_id" validate:"required"`
	ProductID uint      `gorm:"column:product_id" json:"product_id" validate:"required"`
	Product   Product   `gorm:"foreignkey:product_id;association_foreignkey:id"`
}
