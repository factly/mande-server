package model

import (
	"github.com/jinzhu/gorm"
)

// ProductType model
type ProductType struct {
	gorm.Model
	Name string `gorm:"column:name" json:"name" validate:"required"`
}
