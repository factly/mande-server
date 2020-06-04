package model

// Currency model
type Currency struct {
	Base
	IsoCode string `gorm:"column:iso_code" json:"iso_code"`
	Name    string `gorm:"column:name" json:"name"`
}
