package model

import "gorm.io/gorm"

// Currency model
type Currency struct {
	Base
	IsoCode string `gorm:"column:iso_code" json:"iso_code" validate:"required"`
	Name    string `gorm:"column:name" json:"name" validate:"required"`
}

var currencyUser ContextKey = "currency_user"

// BeforeCreate hook
func (currency *Currency) BeforeCreate(tx *gorm.DB) error {
	ctx := tx.Statement.Context
	userID := ctx.Value(currencyUser)

	if userID == nil {
		return nil
	}
	uID := userID.(int)

	currency.CreatedByID = uint(uID)
	currency.UpdatedByID = uint(uID)
	return nil
}
