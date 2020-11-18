package model

// Plan model
type Plan struct {
	Base
	Name        string    `gorm:"column:name" json:"name" validate:"required"`
	Description string    `gorm:"column:description" json:"description"`
	Price       int       `gorm:"column:price" json:"price" validate:"required"`
	CurrencyID  uint      `gorm:"column:currency_id" json:"currency_id" validate:"required"`
	Currency    *Currency `gorm:"foreignKey:currency_id"  json:"currency"`
	Duration    uint      `gorm:"column:duration" json:"duration" validate:"required"`
	Status      string    `gorm:"column:status" json:"status" validate:"required"`
	Catalogs    []Catalog `gorm:"many2many:plan_catalog" json:"catalogs"`
}
