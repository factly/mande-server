package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/factly/data-portal-server/action"
	"github.com/factly/data-portal-server/config"
	"github.com/factly/data-portal-server/model"
	"github.com/factly/data-portal-server/util/meili"
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
	config.SetupVars()

	port, ok := os.LookupEnv("PORT")
	if !ok {
		port = "7720"
	}

	port = ":" + port
	// db setup
	model.SetupDB(config.DSN)

	model.Migration()

	meili.SetupMeiliSearch()

	// register routes
	r := action.RegisterRoutes()

	fmt.Println("swagger-ui http://localhost:7720/swagger/index.html")
	err := http.ListenAndServe(port, r)
	if err != nil {
		log.Fatal(err)
	}
}
