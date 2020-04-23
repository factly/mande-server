package models

import (
	"time"
)

//Status model
type Status struct {
	ID        uint      `gorm:"primary_key"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`
	Name      string    `gorm:"column:name" json:"name" validate:"required"`
}
