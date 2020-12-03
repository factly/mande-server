package model

import "gorm.io/gorm"

// Payment model
type Payment struct {
	Base
	Amount            int      `gorm:"column:amount" json:"amount" validate:"required"`
	Gateway           string   `gorm:"column:gateway" json:"gateway" validate:"required"`
	CurrencyID        uint     `gorm:"column:currency_id" json:"currency_id" validate:"required"`
	Currency          Currency `gorm:"foreignKey:currency_id" json:"currency"`
	Status            string   `gorm:"column:status" json:"status" validate:"required"`
	RazorpayPaymentID string   `gorm:"column:razorpay_payment_id" json:"razorpay_payment_id"`
	RazorpaySignature string   `gorm:"column:razorpay_signature" json:"razorpay_signature"`
}

var paymentUser ContextKey = "payment_user"

// BeforeCreate hook
func (payment *Payment) BeforeCreate(tx *gorm.DB) error {
	ctx := tx.Statement.Context
	userID := ctx.Value(paymentUser)

	if userID == nil {
		return nil
	}
	uID := userID.(int)

	payment.CreatedByID = uint(uID)
	payment.UpdatedByID = uint(uID)
	return nil
}
