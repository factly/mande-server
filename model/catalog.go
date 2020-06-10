package model

import "time"

// Catalog model
type Catalog struct {
	Base
	Title           string    `gorm:"column:title" json:"title"`
	Description     string    `gorm:"column:description" json:"description"`
	Price           int       `gorm:"column:price" json:"price"`
	FeaturedMediaID uint      `gorm:"column:featured_media_id" json:"featured_media_id"`
	PublishedDate   time.Time `gorm:"column:published_date" json:"published_date"`
}

// CatalogProduct model
type CatalogProduct struct {
	Base
	ProductID uint    `gorm:"column:product_id" json:"product_id"`
	Product   Product `gorm:"foreignkey:product_id;association_foreignkey:id"  json:"product"`
	CatalogID uint    `gorm:"column:category_id" json:"category_id"`
}
