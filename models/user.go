package models

import (
	"net/url"
	"time"
)

// User model
type User struct {
	ID        uint      `gorm:"primary_key" json:"id"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`
	Email     string    `gorm:"column:email" json:"email"`
	Name      string    `gorm:"column:name" json:"name"`
	Age       int       `gorm:"column:age" json:"age"`
}

func (u *User) Validate() url.Values {
	errs := url.Values{}

	if u.Age == 0 {
		errs.Add("age", "Age is required")
	}
	if u.Email == "" {
		errs.Add("email", "Email is required")
	}
	if u.Name == "" {
		errs.Add("name", "Name is required")
	}
	return errs
}
