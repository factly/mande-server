package model

// Product model
type Product struct {
	Base
	Title         string      `gorm:"column:title" json:"title" validate:"required"`
	Slug          string      `gorm:"column:slug" json:"slug" validate:"required"`
	Price         int         `gorm:"column:price" json:"price" validate:"required"`
	ProductTypeID uint        `gorm:"column:product_type_id" json:"product_type_id" validate:"required"`
	ProductType   ProductType `gorm:"foreignkey:product_type_id;association_foreignkey:id"`
	StatusID      uint        `gorm:"column:status_id" json:"status_id" validate:"required"`
	Status        Status      `gorm:"foreignkey:status_id;association_foreignkey:id"`
	CurrencyID    uint        `gorm:"column:currency_id" json:"currency_id" validate:"required"`
	Currency      Currency    `gorm:"foreignkey:currency_id;association_foreignkey:id"`
}

// ProductCategory model
type ProductCategory struct {
	Base
	CategoryID uint     `gorm:"column:category_id" json:"category_id" validate:"required"`
	Category   Category `gorm:"foreignkey:category_id;association_foreignkey:id"`
	ProductID  uint     `gorm:"column:product_id" json:"product_id" validate:"required"`
}

// ProductTag model
type ProductTag struct {
	Base
	TagID     uint `gorm:"column:tag_id" json:"tag_id" validate:"required"`
	Tag       Tag  `gorm:"foreignkey:tag_id;association_foreignkey:id"`
	ProductID uint `gorm:"column:product_id" json:"product_id" validate:"required"`
}
