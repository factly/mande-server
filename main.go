package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/factly/data-portal-server/action"
	"github.com/factly/data-portal-server/model"
	"github.com/factly/data-portal-server/validation"
	"github.com/joho/godotenv"
)

// @title Data portal API
// @version 1.0
// @description This is a sample server.

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:7720
// @BasePath /

func main() {
	godotenv.Load()

	port, ok := os.LookupEnv("PORT")
	if !ok {
		port = "7720"
	}

	port = ":" + port
	// db setup
	model.SetupDB()

	// register routes
	r := action.RegisterRoutes()

	// Initialize validator
	validation.InitializeValidator()

	fmt.Println("swagger-ui http://localhost:7720/swagger/index.html")
	http.ListenAndServe(port, r)
}
