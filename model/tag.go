package model

// Tag model
type Tag struct {
	Base
	Title string `gorm:"column:title" json:"title"`
	Slug  string `gorm:"column:slug" json:"slug"`
}
