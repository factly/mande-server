package model

import (
	"github.com/jinzhu/gorm"
)

// Category model
type Category struct {
	gorm.Model
	Title    string `gorm:"column:title" json:"title" validate:"required"`
	Slug     string `gorm:"column:slug" json:"slug" validate:"required"`
	ParentID uint   `gorm:"column:parent_id" json:"parent_id"`
}

/* include meta in category*/
