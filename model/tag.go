package model

import (
	"github.com/jinzhu/gorm"
)

// Tag model
type Tag struct {
	gorm.Model
	Title string `gorm:"column:title" json:"title" validate:"required"`
	Slug  string `gorm:"column:slug" json:"slug" validate:"required"`
}
