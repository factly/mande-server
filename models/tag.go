package models

import (
	"net/url"
	"time"
)

type Tag struct {
	ID        uint      `gorm:"primary_key" json:"id"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`
	Title     string    `gorm:"column:title" json:"title"`
	Slug      string    `gorm:"column:slug" json:"slug"`
}

func (t *Tag) Validate() url.Values {
	errs := url.Values{}

	if t.Title == "" {
		errs.Add("Title", "Title field is required!")
	}
	if t.Slug == "" {
		errs.Add("Slug", "Slug field is required!")
	}
	return errs
}
