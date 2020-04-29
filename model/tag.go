package model

import (
	"time"
)

// Tag model
type Tag struct {
	ID        uint      `gorm:"primary_key" json:"id"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`
	Title     string    `gorm:"column:title" json:"title" validate:"required"`
	Slug      string    `gorm:"column:slug" json:"slug" validate:"required"`
}

/* include meta in tag model*/

// ProductTag model
type ProductTag struct {
	ID        uint      `gorm:"primary_key" json:"id"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`
	TagID     uint      `gorm:"column:tag_id" json:"tag_id" validate:"required"`
	ProductID uint      `gorm:"column:product_id" json:"product_id" validate:"required"`
}
