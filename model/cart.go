package model

// Cart model
type Cart struct {
	Base
	Status string `gorm:"column:status" json:"status"`
	UserID uint   `gorm:"column:user_id" json:"user_id"`
}

// CartItem model
type CartItem struct {
	Base
	CartID    uint    `gorm:"column:cart_id" json:"cart_id"`
	ProductID uint    `gorm:"column:product_id" json:"product_id"`
	Product   Product `gorm:"foreignkey:product_id;association_foreignkey:id" json:"product"`
}
