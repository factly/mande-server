package model

// Order model
type Order struct {
	Base
	UserID    uint    `gorm:"column:user_id" json:"user_id"`
	Status    string  `gorm:"column:status" json:"status"`
	PaymentID uint    `gorm:"column:payment_id" json:"payment_id"`
	Payment   Payment `gorm:"foreignkey:payment_id;association_foreignkey:id" json:"payment"`
	CartID    uint    `gorm:"column:cart_id" json:"cart_id"`
	Cart      Cart    `gorm:"foreignkey:cart_id;association_foreignkey:id" json:"cart"`
}

// OrderItem model
type OrderItem struct {
	Base
	ExtraInfo string  `gorm:"column:extra_info" json:"extra_info"`
	ProductID uint    `gorm:"column:product_id" json:"product_id"`
	Product   Product `gorm:"foreignkey:product_id;association_foreignkey:id"  json:"product"`
	OrderID   uint    `gorm:"column:order_id" json:"order_id"`
}
