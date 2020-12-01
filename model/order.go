package model

import "gorm.io/gorm"

// Order model
type Order struct {
	Base
	UserID          uint      `gorm:"column:user_id" json:"user_id" validate:"required"`
	Status          string    `gorm:"column:status" json:"status" validate:"required"`
	PaymentID       *uint     `gorm:"column:payment_id" json:"payment_id" sql:"DEFAULT:NULL"`
	Payment         *Payment  `gorm:"foreignKey:payment_id" json:"payment"`
	RazorpayOrderID string    `gorm:"column:razorpay_order_id" json:"razorpay_order_id" sql:"DEFAULT:NULL"`
	Products        []Product `gorm:"many2many:order_item;" json:"products"`
}

var orderUser ContextKey = "order_user"

// BeforeCreate hook
func (order *Order) BeforeCreate(tx *gorm.DB) error {
	ctx := tx.Statement.Context
	userID := ctx.Value(orderUser)

	if userID == nil {
		return nil
	}
	uID := userID.(int)

	order.CreatedByID = uint(uID)
	order.UpdatedByID = uint(uID)
	return nil
}
