package model

import "gorm.io/gorm"

// Format model
type Format struct {
	Base
	Name        string `gorm:"column:name" json:"name"`
	Description string `gorm:"column:description" json:"description" `
	IsDefault   bool   `gorm:"column:is_default" json:"is_default" `
}

var formatUser ContextKey = "format_user"

// BeforeCreate hook
func (format *Format) BeforeCreate(tx *gorm.DB) error {
	ctx := tx.Statement.Context
	userID := ctx.Value(formatUser)

	if userID == nil {
		return nil
	}
	uID := userID.(int)

	format.CreatedByID = uint(uID)
	format.UpdatedByID = uint(uID)
	return nil
}
