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
