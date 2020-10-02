package model

// CartItem model
type CartItem struct {
	Base
	Status       string      `gorm:"column:status" json:"status" validate:"required"`
	UserID       uint        `gorm:"column:user_id" json:"user_id" validate:"required"`
	ProductID    uint        `gorm:"column:product_id" json:"product_id" validate:"required"`
	Product      *Product    `gorm:"foreignkey:product_id;association_foreignkey:id" json:"product"`
	MembershipID uint        `gorm:"column:membership_id" json:"membership_id" sql:"DEFAULT:NULL"`
	Membership   *Membership `gorm:"foreignkey:membership_id;association_foreignkey:id" json:"membership"`
}
