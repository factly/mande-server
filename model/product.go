package model

// Product model
type Product struct {
	Base
	Title         string      `gorm:"column:title" json:"title"`
	Slug          string      `gorm:"column:slug" json:"slug"`
	Price         int         `gorm:"column:price" json:"price"`
	ProductTypeID uint        `gorm:"column:product_type_id" json:"product_type_id"`
	ProductType   ProductType `gorm:"foreignkey:product_type_id;association_foreignkey:id" json:"product_type"`
	Status        string      `gorm:"column:status" json:"status"`
	CurrencyID    uint        `gorm:"column:currency_id" json:"currency_id"`
	Currency      Currency    `gorm:"foreignkey:currency_id;association_foreignkey:id"  json:"currency"`
}

// ProductCategory model
type ProductCategory struct {
	Base
	CategoryID uint     `gorm:"column:category_id" json:"category_id"`
	Category   Category `gorm:"foreignkey:category_id;association_foreignkey:id"  json:"category"`
	ProductID  uint     `gorm:"column:product_id" json:"product_id"`
}

// ProductTag model
type ProductTag struct {
	Base
	TagID     uint `gorm:"column:tag_id" json:"tag_id"`
	Tag       Tag  `gorm:"foreignkey:tag_id;association_foreignkey:id"  json:"tag"`
	ProductID uint `gorm:"column:product_id" json:"product_id"`
}
