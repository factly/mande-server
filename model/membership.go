package model

import "gorm.io/gorm"

// Membership model
type Membership struct {
	Base
	Status          string   `gorm:"column:status" json:"status" validate:"required"`
	UserID          uint     `gorm:"column:user_id" json:"user_id" validate:"required"`
	OrganisationID  uint     `gorm:"column:organisation_id" json:"organisation_id" validate:"required"`
	PaymentID       *uint    `gorm:"column:payment_id" json:"payment_id" sql:"DEFAULT:NULL"`
	Payment         *Payment `gorm:"foreignKey:payment_id" json:"payment"`
	PlanID          uint     `gorm:"column:plan_id" json:"plan_id" validate:"required"`
	Plan            Plan     `gorm:"foreignKey:plan_id" json:"plan"`
	RazorpayOrderID string   `gorm:"column:razorpay_order_id" json:"razorpay_order_id" sql:"DEFAULT:NULL"`
}

var membershipUser ContextKey = "membership_user"

// BeforeCreate hook
func (membership *Membership) BeforeCreate(tx *gorm.DB) error {
	ctx := tx.Statement.Context
	userID := ctx.Value(membershipUser)

	if userID == nil {
		return nil
	}
	uID := userID.(int)

	membership.CreatedByID = uint(uID)
	membership.UpdatedByID = uint(uID)
	return nil
}
