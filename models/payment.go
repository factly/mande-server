package models

import (
	"net/url"
	"time"
)

// Payment model
type Payment struct {
	ID         uint      `gorm:"primary_key" json:"id"`
	CreatedAt  time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt  time.Time `gorm:"column:updated_at" json:"updated_at"`
	Amount     int       `gorm:"column:amount" json:"amount"`
	Gateway    string    `gorm:"column:gateway" json:"gateway"`
	CurrencyID uint      `gorm:"column:currency_id" json:"currency_id"`
	Currency   Currency  `gorm:"foreignkey:currency_id;association_foreignkey:id"`
	Status     string    `gorm:"column:status" json:"status"`
}

// validation logic
func (p *Payment) Validate() url.Values {
	errs := url.Values{}

	if p.Status == "" {
		errs.Add("status", "status field is required!")
	}

	if p.Amount == 0 {
		errs.Add("amount", "amount is required!")
	}

	if p.CurrencyID == 0 {
		errs.Add("currencyid", "Currency id is required!")
	}

	if p.Gateway == "" {
		errs.Add("gateway", "Gateway is required!")
	}

	return errs
}
