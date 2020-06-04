package model

// Plan model
type Plan struct {
	Base
	PlanName string `gorm:"column:plan_name" json:"plan_name"`
	PlanInfo string `gorm:"column:plan_info" json:"plan_info"`
	Status   string `gorm:"column:status" json:"status"`
}
