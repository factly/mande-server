package model

//Status model
type Status struct {
	BaseModel
	Name string `gorm:"column:name" json:"name" validate:"required"`
}
