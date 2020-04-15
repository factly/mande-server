package models

import (
	"net/url"
	"time"
)

// Plan model
type Plan struct {
	ID        uint      `gorm:"primary_key"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`
	PlanName  string    `gorm:"column:plan_name" json:"plan_name"`
	PlanInfo  string    `gorm:"column:plan_info" json:"plan_info"`
	Status    string    `gorm:"column:status" json:"status"`
}

func (p *Plan) Validate() url.Values {
	errs := url.Values{}

	if p.PlanName == "" {
		errs.Add("plan_name", "Plan Name is required")
	}
	if p.PlanInfo == "" {
		errs.Add("plan_info", "Plan Info is required")
	}
	if p.Status == "" {
		errs.Add("status", "Status is required")
	}

	return errs
}
