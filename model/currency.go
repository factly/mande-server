package model

import (
	"github.com/jinzhu/gorm"
)

// Currency model
type Currency struct {
	gorm.Model
	IsoCode string `gorm:"column:iso_code" json:"iso_code" validate:"required"`
	Name    string `gorm:"column:name" json:"name" validate:"required"`
}
