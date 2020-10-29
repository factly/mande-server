package model

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
