package model

// Membership model
type Membership struct {
	Base
	Status          string   `gorm:"column:status" json:"status" validate:"required"`
	UserID          uint     `gorm:"column:user_id" json:"user_id" validate:"required"`
	PaymentID       uint     `gorm:"column:payment_id" json:"payment_id" sql:"DEFAULT:NULL"`
	Payment         *Payment `gorm:"foreignkey:payment_id;association_foreignkey:id" json:"payment"`
	PlanID          uint     `gorm:"column:plan_id" json:"plan_id" validate:"required"`
	Plan            Plan     `gorm:"foreignkey:plan_id;association_foreignkey:id" json:"plan"`
	RazorpayOrderID string   `gorm:"column:razorpay_order_id" json:"razorpay_order_id" sql:"DEFAULT:NULL"`
}
