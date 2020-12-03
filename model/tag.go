package model

import "gorm.io/gorm"

// Tag model
type Tag struct {
	Base
	Title    string    `gorm:"column:title" json:"title" validate:"required"`
	Slug     string    `gorm:"column:slug" json:"slug" validate:"required"`
	Products []Product `gorm:"many2many:product_tag;" json:"products"`
	Datasets []Dataset `gorm:"many2many:dataset_tag;" json:"datasets"`
}

var tagUser ContextKey = "tag_user"

// BeforeCreate hook
func (tag *Tag) BeforeCreate(tx *gorm.DB) error {
	ctx := tx.Statement.Context
	userID := ctx.Value(tagUser)

	if userID == nil {
		return nil
	}
	uID := userID.(int)

	tag.CreatedByID = uint(uID)
	tag.UpdatedByID = uint(uID)
	return nil
}
