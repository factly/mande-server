package model

// User model
type User struct {
	Base
	Email     string `gorm:"column:email" json:"email"`
	FirstName string `gorm:"column:first_name" json:"first_name"`
	LastName  string `gorm:"column:last_name" json:"last_name"`
}
