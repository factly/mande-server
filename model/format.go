package model

// Format model
type Format struct {
	Base
	Name        string `gorm:"column:name" json:"name"`
	Description string `gorm:"column:description" json:"description" `
	IsDefault   bool   `gorm:"column:is_default" json:"is_default" `
}
