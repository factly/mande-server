package model

// Category model
type Category struct {
	Base
	Title    string `gorm:"column:title" json:"title" validate:"required"`
	Slug     string `gorm:"column:slug" json:"slug" validate:"required"`
	ParentID uint   `gorm:"column:parent_id" json:"parent_id" validate:"required"`
}

/* include meta in category*/
