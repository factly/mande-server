package model

// CartItem model
type CartItem struct {
	Base
	Status    string   `gorm:"column:status" json:"status" validate:"required"`
	UserID    uint     `gorm:"column:user_id" json:"user_id" validate:"required"`
	ProductID uint     `gorm:"column:product_id" json:"product_id" validate:"required"`
	Product   *Product `gorm:"foreignkey:product_id;association_foreignkey:id" json:"product"`
}
