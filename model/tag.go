package model

// Tag model
type Tag struct {
	BaseModel
	Title string `gorm:"column:title" json:"title" validate:"required"`
	Slug  string `gorm:"column:slug" json:"slug" validate:"required"`
}
