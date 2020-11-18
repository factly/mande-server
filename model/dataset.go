package model

import "github.com/jinzhu/gorm/dialects/postgres"

// Dataset model
type Dataset struct {
	Base
	Title            string         `gorm:"column:title" json:"title"`
	Description      string         `gorm:"column:description" json:"description"`
	Source           string         `gorm:"column:source" json:"source"`
	Frequency        string         `gorm:"column:frequency" json:"frequency"`
	TemporalCoverage string         `gorm:"column:temporal_coverage" json:"temporal_coverage"`
	Granularity      string         `gorm:"column:granularity" json:"granularity"`
	ContactName      string         `gorm:"column:contact_name" json:"contact_name"`
	ContactEmail     string         `gorm:"column:contact_email" json:"contact_email"`
	License          string         `gorm:"column:license" json:"license"`
	DataStandard     string         `gorm:"column:data_standard" json:"data_standard"`
	SampleURL        string         `gorm:"column:sample_url" json:"sample_url"`
	RelatedArticles  postgres.Jsonb `gorm:"column:related_articles" json:"related_articles" swaggertype:"primitive,string"`
	TimeSaved        int            `gorm:"column:time_saved" json:"time_saved"`
	Price            int            `gorm:"column:price" json:"price" validate:"required"`
	CurrencyID       uint           `gorm:"column:currency_id" json:"currency_id" validate:"required"`
	Currency         *Currency      `gorm:"foreignkey:currency_id" json:"currency"`
	FeaturedMediumID *uint          `gorm:"column:featured_medium_id;default:NULL" json:"featured_medium_id"`
	FeaturedMedium   *Medium        `gorm:"foreignKey:featured_medium_id" json:"featured_medium"`
	Products         []Product      `gorm:"many2many:product_dataset;" json:"products"`
	Tags             []Tag          `gorm:"many2many:dataset_tag;" json:"tags"`
}

// DatasetFormat model
type DatasetFormat struct {
	Base
	FormatID  uint   `gorm:"column:format_id" json:"format_id"`
	Format    Format `gorm:"foreignKey:format_id"  json:"format"`
	DatasetID uint   `gorm:"column:dataset_id" json:"dataset_id"`
	URL       string `gorm:"column:url" json:"url"`
}
