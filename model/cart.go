package model

// Cart model
type Cart struct {
	Base
	Status   string    `gorm:"column:status" json:"status" validate:"required"`
	UserID   uint      `gorm:"column:user_id" json:"user_id" validate:"required"`
	Products []Product `gorm:"many2many:cart_item;" json:"products"`
}
