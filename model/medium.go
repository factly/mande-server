package model

// Medium model
type Medium struct {
	Base
	Name        string `gorm:"column:name" json:"name"`
	Slug        string `gorm:"column:slug" json:"slug"`
	Type        string `gorm:"column:type" json:"type"`
	Title       string `gorm:"column:title" json:"title"`
	Description string `gorm:"column:description" json:"description"`
	Caption     string `gorm:"column:caption" json:"caption"`
	AltText     string `gorm:"column:alt_text" json:"alt_text"`
	FileSize    int    `gorm:"column:file_size" json:"file_size"`
	URL         string `gorm:"column:url" json:"url"`
	Dimensions  string `gorm:"column:dimensions" json:"dimensions"`
}
