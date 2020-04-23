package models

import (
	"time"
)

// Membership model
type Membership struct {
	ID        uint      `gorm:"primary_key" json:"id"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`
	Status    string    `gorm:"column:status" json:"status" validate:"required"`
	UserID    uint      `gorm:"column:user_id" json:"user_id" validate:"required"`
	User      User      `gorm:"foreignkey:user_id;association_foreignkey:id"`
	PaymentID uint      `gorm:"column:payment_id" json:"payment_id" validate:"required"`
	Payment   Payment   `gorm:"foreignkey:payment_id;association_foreignkey:id"`
	PlanID    uint      `gorm:"column:plan_id" json:"plan_id" validate:"required"`
	Plan      Plan      `gorm:"foreignkey:plan_id;association_foreignkey:id"`
}
