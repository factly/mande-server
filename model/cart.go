package model

import (
	"github.com/jinzhu/gorm"
)

// Cart model
type Cart struct {
	gorm.Model
	Status string `gorm:"column:status" json:"status" validate:"required"`
	UserID uint   `gorm:"column:user_id" json:"user_id" validate:"required"`
}

// CartItem model
type CartItem struct {
	gorm.Model
	IsDeleted bool    `gorm:"column:is_deleted" json:"is_deleted"`
	CartID    uint    `gorm:"column:cart_id" json:"cart_id" validate:"required"`
	ProductID uint    `gorm:"column:product_id" json:"product_id" validate:"required"`
	Product   Product `gorm:"foreignkey:product_id;association_foreignkey:id"`
}
