package model

import (
	"github.com/jinzhu/gorm"
)

//Status model
type Status struct {
	gorm.Model
	Name string `gorm:"column:name" json:"name" validate:"required"`
}
