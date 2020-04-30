package model

import (
	"github.com/jinzhu/gorm"
)

// Plan model
type Plan struct {
	gorm.Model
	PlanName string `gorm:"column:plan_name" json:"plan_name" validate:"required"`
	PlanInfo string `gorm:"column:plan_info" json:"plan_info" validate:"required"`
	Status   string `gorm:"column:status" json:"status" validate:"required"`
}
