package models

import (
	"net/url"
	"time"
)

//Status model
type Status struct {
	ID        uint      `gorm:"primary_key"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`
	Name      string    `gorm:"column:name" json:"name"`
}

func (s *Status) Validate() url.Values {
	errs := url.Values{}

	if s.Name == "" {
		errs.Add("Name", "name field is required!")
	}

	return errs
}
