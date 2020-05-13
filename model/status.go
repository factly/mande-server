package model

//Status model
type Status struct {
	Base
	Name string `gorm:"column:name" json:"name" validate:"required"`
}
