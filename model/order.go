package model

import (
	"github.com/jinzhu/gorm"
)

// Order model
type Order struct {
	gorm.Model
	UserID    uint    `gorm:"column:user_id" json:"user_id" validate:"required"`
	Status    string  `gorm:"column:status" json:"status" validate:"required"`
	PaymentID uint    `gorm:"column:payment_id" json:"payment_id" validate:"required"`
	Payment   Payment `gorm:"foreignkey:payment_id;association_foreignkey:id"`
	CartID    uint    `gorm:"column:cart_id" json:"cart_id" validate:"required"`
	Cart      Cart    `gorm:"foreignkey:cart_id;association_foreignkey:id"`
}

// OrderItem model
type OrderItem struct {
	gorm.Model
	ExtraInfo string  `gorm:"column:extra_info" json:"extra_info" validate:"required"`
	ProductID uint    `gorm:"column:product_id" json:"product_id" validate:"required"`
	Product   Product `gorm:"foreignkey:product_id;association_foreignkey:id"`
	OrderID   uint    `gorm:"column:order_id" json:"order_id" validate:"required"`
	Order     Order   `gorm:"foreignkey:order_id;association_foreignkey:id"`
}
