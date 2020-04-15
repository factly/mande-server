package models

import "net/url"

// User model
type User struct {
	ID    uint   `gorm:"primary_key" json:"id"`
	Email string `gorm:"column:email" json:"email"`
	Name  string `gorm:"column:name" json:"name"`
	Age   int    `gorm:"column:age" json:"age"`
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
