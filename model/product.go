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
}

// ProductDataset model
type ProductDataset struct {
	Base
	DatasetID uint    `gorm:"column:dataset_id" json:"dataset_id" validate:"required"`
	Dataset   Dataset `gorm:"foreignkey:dataset_id;association_foreignkey:id"  json:"dataset"`
	ProductID uint    `gorm:"column:product_id" json:"product_id" validate:"required"`
}

// ProductTag model
type ProductTag struct {
	Base
	TagID     uint `gorm:"column:tag_id" json:"tag_id" validate:"required"`
	Tag       Tag  `gorm:"foreignkey:tag_id;association_foreignkey:id"  json:"tag"`
	ProductID uint `gorm:"column:product_id" json:"product_id" validate:"required"`
}
