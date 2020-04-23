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

	// register routes
	r := actions.RegisterRoutes()

	fmt.Println("swagger-ui http://localhost:3000/swagger/index.html")
	http.ListenAndServe(":3000", r)
}
