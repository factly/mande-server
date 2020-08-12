package model

import "time"

// Catalog model
type Catalog struct {
	Base
	Title            string    `gorm:"column:title" json:"title"`
	Description      string    `gorm:"column:description" json:"description"`
	Price            int       `gorm:"column:price" json:"price"`
	CurrencyID       uint      `gorm:"column:currency_id" json:"currency_id" validate:"required"`
	Currency         *Currency `gorm:"foreignkey:currency_id;association_foreignkey:id"  json:"currency"`
	FeaturedMediumID uint      `gorm:"column:featured_medium_id" json:"featured_medium_id" sql:"DEFAULT:NULL"`
	FeaturedMedium   *Medium   `gorm:"foreignkey:featured_medium_id;association_foreignkey:id"  json:"featured_medium"`
	PublishedDate    time.Time `gorm:"column:published_date" json:"published_date"`
	Products         []Product `gorm:"many2many:catalog_product;" json:"products"`
}
