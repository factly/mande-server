package model

import "gorm.io/gorm"

// Plan model
type Plan struct {
	Base
	Name        string    `gorm:"column:name" json:"name" validate:"required"`
	Description string    `gorm:"column:description" json:"description"`
	Price       int       `gorm:"column:price" json:"price" validate:"required"`
	Users       int       `gorm:"column:users" json:"users" validate:"required"`
	CurrencyID  uint      `gorm:"column:currency_id" json:"currency_id" validate:"required"`
	Currency    *Currency `gorm:"foreignKey:currency_id"  json:"currency"`
	Duration    uint      `gorm:"column:duration" json:"duration" validate:"required"`
	Status      string    `gorm:"column:status" json:"status" validate:"required"`
	AllProducts bool      `gorm:"column:all_products" json:"all_products"`
	Catalogs    []Catalog `gorm:"many2many:plan_catalog" json:"catalogs"`
}

var planUser ContextKey = "plan_user"

// BeforeCreate hook
func (plan *Plan) BeforeCreate(tx *gorm.DB) error {
	ctx := tx.Statement.Context
	userID := ctx.Value(planUser)

	if userID == nil {
		return nil
	}
	uID := userID.(int)

	plan.CreatedByID = uint(uID)
	plan.UpdatedByID = uint(uID)
	return nil
}
