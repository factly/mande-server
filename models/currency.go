package models

import (
	"net/url"
	"time"
)

// Currency model
type Currency struct {
	ID        uint      `gorm:"primary_key"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`
	IsoCode   string    `gorm:"column:iso_code" json:"iso_code"`
	Name      string    `gorm:"column:name" json:"name"`
}

func (p *Currency) Validate() url.Values {
	errs := url.Values{}

	if p.IsoCode == "" {
		errs.Add("Iso code", "Iso code is required")
	}
	if p.Name == "" {
		errs.Add("Name", "Name is required")
	}

	return errs
}
