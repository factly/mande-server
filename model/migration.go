package model

// Migration does database migrations
func Migration() {
	// db migrations
	DB.AutoMigrate(
		&Currency{},
		&Payment{},
		&Membership{},
		&Plan{},
		&Product{},
		&Tag{},
		&Catalog{},
		&CartItem{},
		&Order{},
		&Dataset{},
		&Format{},
		&DatasetFormat{},
		&Medium{},
	)

	// Adding foreignKey
	DB.Model(&Payment{}).AddForeignKey("currency_id", "dp_currency(id)", "RESTRICT", "RESTRICT")
	DB.Model(&Membership{}).AddForeignKey("plan_id", "dp_plan(id)", "RESTRICT", "RESTRICT")
	DB.Model(&Membership{}).AddForeignKey("payment_id", "dp_payment(id)", "RESTRICT", "RESTRICT")
	DB.Model(&Product{}).AddForeignKey("currency_id", "dp_currency(id)", "RESTRICT", "RESTRICT")
	DB.Model(&Product{}).AddForeignKey("featured_medium_id", "dp_medium(id)", "RESTRICT", "RESTRICT")
	DB.Model(&CartItem{}).AddForeignKey("product_id", "dp_product(id)", "RESTRICT", "RESTRICT")
	DB.Model(&CartItem{}).AddForeignKey("membership_id", "dp_membership(id)", "RESTRICT", "RESTRICT")
	DB.Model(&Order{}).AddForeignKey("payment_id", "dp_payment(id)", "RESTRICT", "RESTRICT")
	DB.Model(&DatasetFormat{}).AddForeignKey("format_id", "dp_format(id)", "RESTRICT", "RESTRICT")
	DB.Model(&Catalog{}).AddForeignKey("featured_medium_id", "dp_medium(id)", "RESTRICT", "RESTRICT")
	DB.Model(&Dataset{}).AddForeignKey("featured_medium_id", "dp_medium(id)", "RESTRICT", "RESTRICT")
	DB.Model(&Dataset{}).AddForeignKey("currency_id", "dp_currency(id)", "RESTRICT", "RESTRICT")
}
