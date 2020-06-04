package model

// Category model
type Category struct {
	Base
	Title    string `gorm:"column:title" json:"title"`
	Slug     string `gorm:"column:slug" json:"slug"`
	ParentID uint   `gorm:"column:parent_id" json:"parent_id"`
}

/* include meta in category*/
