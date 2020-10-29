package model

// Migration does database migrations
func Migration() {
	// db migrations
	_ = DB.AutoMigrate(
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
}
