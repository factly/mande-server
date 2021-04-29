package model

import (
	"time"

	"gorm.io/gorm"
)

// Catalog model
type Catalog struct {
	Base
	Title            string     `gorm:"column:title" json:"title"`
	Slug             string     `gorm:"column:slug" json:"slug"`
	Price            int        `gorm:"column:price" json:"price"`
	CurrencyID       uint       `gorm:"column:currency_id" json:"currency_id"`
	Currency         *Currency  `gorm:"foreignKey:currency_id"  json:"currency"`
	Description      string     `gorm:"column:description" json:"description"`
	FeaturedMediumID *uint      `gorm:"column:featured_medium_id;default:NULL" json:"featured_medium_id"`
	FeaturedMedium   *Medium    `gorm:"foreignKey:featured_medium_id"  json:"featured_medium"`
	PublishedDate    *time.Time `gorm:"column:published_date" json:"published_date" sql:"DEFAULT:NULL"`
	Plans            []Plan     `gorm:"many2many:plan_catalog;" json:"plans"`
	Products         []Product  `gorm:"many2many:catalog_product;" json:"products"`
}

var catalogUser ContextKey = "catalog_user"

// BeforeCreate hook
func (catalog *Catalog) BeforeCreate(tx *gorm.DB) error {
	ctx := tx.Statement.Context
	userID := ctx.Value(catalogUser)

	if userID == nil {
		return nil
	}
	uID := userID.(int)

	catalog.CreatedByID = uint(uID)
	catalog.UpdatedByID = uint(uID)
	return nil
}
