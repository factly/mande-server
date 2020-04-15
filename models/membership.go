package models

import (
	"net/url"
	"time"
)

// Membership model
type Membership struct {
	ID        uint      `gorm:"primary_key"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`
	Status    string    `gorm:"column:status" json:"status"`
	UserID    uint      `gorm:"column:user_id" json:"user_id"`
	User      User      `gorm:"foreignkey:user_id;association_foreignkey:id"`
	PaymentID uint      `gorm:"column:payment_id" json:"payment_id"`
	Payment   Payment   `gorm:"foreignkey:payment_id;association_foreignkey:id"`
	PlanID    uint      `gorm:"column:plan_id" json:"plan_id"`
	Plan      Plan      `gorm:"foreignkey:plan_id;association_foreignkey:id"`
}

// validation logic
func (p *Membership) Validate() url.Values {
	errs := url.Values{}

	if p.Status == "" {
		errs.Add("status", "status field is required!")
	}

	if p.UserID == 0 {
		errs.Add("user_id", "user id is required!")
	}

	if p.PaymentID == 0 {
		errs.Add("payment_id", "payment id is required!")
	}

	if p.PlanID == 0 {
		errs.Add("plan_id", "plan id is required!")
	}

	return errs
}
