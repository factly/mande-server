package model

// Tag model
type Tag struct {
	Base
	Title    string    `gorm:"column:title" json:"title" validate:"required"`
	Slug     string    `gorm:"column:slug" json:"slug" validate:"required"`
	Products []Product `gorm:"many2many:product_tag;" json:"products"`
}
