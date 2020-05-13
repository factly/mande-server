package model

// ProductType model
type ProductType struct {
	Base
	Name string `gorm:"column:name" json:"name" validate:"required"`
}
