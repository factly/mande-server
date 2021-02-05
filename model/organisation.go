package model

// OrgWithRole model definition
type OrgWithRole struct {
	Organisation
	Permission OrganisationUser `json:"permission"`
}

// Organisation model definition
type Organisation struct {
	Base
	Title            string  `json:"title"`
	Slug             string  `json:"slug"`
	Description      string  `json:"description"`
	FeaturedMediumID *uint   `json:"featured_medium_id"`
	Medium           *Medium `json:"medium"`
}

// OrganisationUser model definition
type OrganisationUser struct {
	Base
	UserID         uint   `json:"user_id"`
	OrganisationID uint   `json:"organisation_id"`
	Role           string `json:"role"`
}

// User model definition
type User struct {
	Base
	Email     string `gorm:"column:email;unique_index" json:"email"`
	FirstName string `gorm:"column:first_name" json:"first_name"`
	LastName  string `gorm:"column:last_name" json:"last_name"`
	BirthDate string `gorm:"column:birth_date" json:"birth_date"`
	Gender    string `gorm:"column:gender" json:"gender"`
}
