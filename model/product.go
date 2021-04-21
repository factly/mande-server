package model

import "gorm.io/gorm"

// Product model
type Product struct {
	Base
	Title            string    `gorm:"column:title" json:"title" validate:"required"`
	Description      string    `gorm:"column:description" json:"description"`
	Slug             string    `gorm:"column:slug" json:"slug" validate:"required"`
	Price            int       `gorm:"column:price" json:"price" validate:"required"`
	Status           string    `gorm:"column:status" json:"status" validate:"required"`
	FeaturedMediumID *uint     `gorm:"column:featured_medium_id;default:NULL" json:"featured_medium_id"`
	FeaturedMedium   *Medium   `gorm:"foreignKey:featured_medium_id"  json:"featured_medium"`
	CurrencyID       uint      `gorm:"column:currency_id" json:"currency_id" validate:"required"`
	Currency         *Currency `gorm:"foreignKey:currency_id"  json:"currency"`
	Catalogs         []Catalog `gorm:"many2many:catalog_product;" json:"catalogs"`
	Tags             []Tag     `gorm:"many2many:product_tag;" json:"tags"`
	Datasets         []Dataset `gorm:"many2many:product_dataset;" json:"datasets"`
	Orders           []Order   `gorm:"many2many:order_item;" json:"orders"`
}

var productUser ContextKey = "product_user"

// BeforeCreate hook
func (product *Product) BeforeCreate(tx *gorm.DB) error {
	ctx := tx.Statement.Context
	userID := ctx.Value(productUser)

	if userID == nil {
		return nil
	}
	uID := userID.(int)

	product.CreatedByID = uint(uID)
	product.UpdatedByID = uint(uID)
	return nil
}
