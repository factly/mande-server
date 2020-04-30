package model

import (
	"github.com/jinzhu/gorm"
)

// User model
type User struct {
	gorm.Model
	Email     string `gorm:"column:email" json:"email" validate:"required"`
	FirstName string `gorm:"column:first_name" json:"first_name" validate:"required"`
	LastName  string `gorm:"column:last_name" json:"last_name" validate:"required"`
}
