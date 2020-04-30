package model

import (
	"time"
)

// Category model
type Category struct {
	ID        uint      `gorm:"primary_key" json:"id"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`
	Title     string    `gorm:"column:title" json:"title" validate:"required"`
	Slug      string    `gorm:"column:slug" json:"slug" validate:"required"`
	ParentID  uint      `gorm:"column:parent_id" json:"parent_id" validate:"required"`
}

/* include meta in category*/
