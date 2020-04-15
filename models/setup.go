package models

import (
	"log"
	"os"
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/joho/godotenv"
)

var DB *gorm.DB

// SetupDB is database setuo
func SetupDB() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("error loding .env file")
	}
	fmt.Println("connecting to database ...")
	username := os.Getenv("db_user")
	password := os.Getenv("db_pass")
	dbName := os.Getenv("db_name")
	dbHost := os.Getenv("db_host")

	connStr := fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable password=%s", dbHost, username, dbName, password) //Build connection string
	//connStr := "user=postgres dbname=data_portal host=localhost sslmode=disable password=postgres"

	DB, err = gorm.Open("postgres", connStr)

	if err != nil {
		log.Fatal(err)
	}

}
