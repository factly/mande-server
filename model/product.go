package model

// Product model
type Product struct {
	Base
	Title            string    `gorm:"column:title" json:"title" validate:"required"`
	Slug             string    `gorm:"column:slug" json:"slug" validate:"required"`
	Price            int       `gorm:"column:price" json:"price" validate:"required"`
	Status           string    `gorm:"column:status" json:"status" validate:"required"`
	FeaturedMediumID uint      `gorm:"column:featured_medium_id" json:"featured_medium_id" sql:"DEFAULT:NULL"`
	FeaturedMedium   *Medium   `gorm:"foreignkey:featured_medium_id;association_foreignkey:id"  json:"featured_medium"`
	CurrencyID       uint      `gorm:"column:currency_id" json:"currency_id" validate:"required"`
	Currency         *Currency `gorm:"foreignkey:currency_id;association_foreignkey:id"  json:"currency"`
	Catalogs         []Catalog `gorm:"many2many:catalog_product;" json:"catalogs"`
	Tags             []Tag     `gorm:"many2many:product_tag;" json:"tags"`
	Datasets         []Dataset `gorm:"many2many:product_dataset;" json:"datasets"`
	Carts            []Cart    `gorm:"many2many:cart_item;" json:"carts"`
}
