package model

import "gorm.io/gorm"

// CartItem model
type CartItem struct {
	Base
	Status       string      `gorm:"column:status" json:"status" validate:"required"`
	UserID       uint        `gorm:"column:user_id" json:"user_id" validate:"required"`
	ProductID    uint        `gorm:"column:product_id" json:"product_id" validate:"required"`
	Product      *Product    `gorm:"foreignKey:product_id" json:"product"`
	MembershipID *uint       `gorm:"column:membership_id;default:NULL" json:"membership_id"`
	Membership   *Membership `gorm:"foreignKey:membership_id" json:"membership"`
}

var cartitemUser ContextKey = "cartitem_user"

// BeforeCreate hook
func (cartitem *CartItem) BeforeCreate(tx *gorm.DB) error {
	ctx := tx.Statement.Context
	userID := ctx.Value(cartitemUser)

	if userID == nil {
		return nil
	}
	uID := userID.(int)

	cartitem.CreatedByID = uint(uID)
	cartitem.UpdatedByID = uint(uID)
	return nil
}
