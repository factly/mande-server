package models

import (
	"net/url"
	"time"
)

// Poduct model
type Product struct {
	ID            uint        `gorm:"primary_key"`
	CreatedAt     time.Time   `gorm:"column:created_at" json:"created_at"`
	UpdatedAt     time.Time   `gorm:"column:updated_at" json:"updated_at"`
	Title         string      `gorm:"column:title" json:"title"`
	Slug          string      `gorm:"column:slug" json:"slug"`
	Price         int         `gorm:"column:price" json:"price"`
	ProductTypeID uint        `gorm:"column:product_type_id" json:"product_type_id"`
	ProductType   ProductType `gorm:"foreignkey:product_type_id;association_foreignkey:id"`
	StatusID      uint        `gorm:"column:status_id" json:"status_id"`
	Status        Status      `gorm:"foreignkey:status_id;association_foreignkey:id"`
	CurrencyID    uint        `gorm:"column:currency_id" json:"currency_id"`
	Currency      Currency    `gorm:"foreignkey:currency_id;association_foreignkey:id"`
}

func (p *Product) Validate() url.Values {
	errs := url.Values{}

	if p.Title == "" {
		errs.Add("Title", "Title field is required!")
	}
	if p.Slug == "" {
		errs.Add("Slug", "Slug field is required!")
	}
	if p.CurrencyID == 0 {
		errs.Add("Currency Id", "Currency id is required!")
	}

	return errs
}
