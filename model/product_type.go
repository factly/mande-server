package model

// ProductType model
type ProductType struct {
	BaseModel
	Name string `gorm:"column:name" json:"name" validate:"required"`
}
