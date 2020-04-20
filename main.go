package main

import (
	"fmt"
	"net/http"

	"github.com/factly/data-portal-api/models"

	"github.com/factly/data-portal-api/actions"
)

// @title Data portal API
// @version 1.0
// @description This is a sample server.

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:3000
// @BasePath /

func main() {
	// db setup
	models.SetupDB()

	// db migrations
	models.DB.AutoMigrate(
		&models.Currency{},
		&models.Payment{},
		&models.Membership{},
		&models.Plan{},
		&models.User{},
		&models.Product{},
		&models.ProductType{},
		&models.Status{},
		&models.Tag{},
	)
	// Adding foreignKey
	models.DB.Model(&models.Payment{}).AddForeignKey("currency_id", "currencies(id)", "RESTRICT", "RESTRICT")
	models.DB.Model(&models.Membership{}).AddForeignKey("user_id", "users(id)", "RESTRICT", "RESTRICT")
	models.DB.Model(&models.Membership{}).AddForeignKey("plan_id", "plans(id)", "RESTRICT", "RESTRICT")
	models.DB.Model(&models.Membership{}).AddForeignKey("payment_id", "payments(id)", "RESTRICT", "RESTRICT")
	models.DB.Model(&models.Product{}).AddForeignKey("currency_id", "currencies(id)", "RESTRICT", "RESTRICT")
	models.DB.Model(&models.Product{}).AddForeignKey("status_id", "statuses(id)", "RESTRICT", "RESTRICT")
	models.DB.Model(&models.Product{}).AddForeignKey("product_type_id", "product_types(id)", "RESTRICT", "RESTRICT")

	r := actions.RegisterRoutes()
	fmt.Println("swagger-ui http://localhost:3000/swagger/index.html")
	http.ListenAndServe(":3000", r)
}
