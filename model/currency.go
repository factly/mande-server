package model

// Currency model
type Currency struct {
	Base
	IsoCode string `gorm:"column:iso_code" json:"iso_code" validate:"required"`
	Name    string `gorm:"column:name" json:"name" validate:"required"`
}
