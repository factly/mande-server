package model

// Plan model
type Plan struct {
	Base
	Name        string    `gorm:"column:name" json:"plan_name" validate:"required"`
	Description string    `gorm:"column:description" json:"description"`
	Duration    uint      `gorm:"column:duration" json:"duration" validate:"required"`
	Status      string    `gorm:"column:status" json:"status" validate:"required"`
	Catalogs    []Catalog `gorm:"many2many:plan_catalog" json:"catalogs"`
}
