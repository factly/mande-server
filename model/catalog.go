package model

import "time"

// Catalog model
type Catalog struct {
	Base
	Title            string    `gorm:"column:title" json:"title"`
	Description      string    `gorm:"column:description" json:"description"`
	FeaturedMediumID uint      `gorm:"column:featured_medium_id" json:"featured_medium_id" sql:"DEFAULT:NULL"`
	FeaturedMedium   *Medium   `gorm:"foreignKey:featured_medium_id"  json:"featured_medium"`
	PublishedDate    time.Time `gorm:"column:published_date" json:"published_date"`
	Plans            []Plan    `gorm:"many2many:plan_catalog;" json:"plans"`
	Products         []Product `gorm:"many2many:catalog_product;" json:"products"`
}
