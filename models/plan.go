package models

import (
	"time"
)

// Plan model
type Plan struct {
	ID        uint      `gorm:"primary_key"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`
	PlanName  string    `gorm:"column:plan_name" json:"plan_name" validate:"required"`
	PlanInfo  string    `gorm:"column:plan_info" json:"plan_info" validate:"required"`
	Status    string    `gorm:"column:status" json:"status" validate:"required"`
}
