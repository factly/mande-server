package models

import (
	"fmt"
	"log"
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres" // postgres
	"github.com/joho/godotenv"
)

// DB - gorm DB
var DB *gorm.DB

// SetupDB is database setuo
func SetupDB() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("error loding .env file")
	}
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbHost := os.Getenv("DB_HOST")

	connStr := fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable password=%s", dbHost, dbUser, dbName, dbPassword) //Build connection string
	//connStr := "user=postgres dbname=data_portal host=localhost sslmode=disable password=postgres"

	DB, err = gorm.Open("postgres", connStr)

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("connected to database ...")

	// db migrations
	DB.AutoMigrate(
		&Currency{},
		&Payment{},
		&Membership{},
		&Plan{},
		&User{},
		&Product{},
		&ProductType{},
		&Status{},
		&Tag{},
		&ProductTag{},
		&Category{},
		&ProductCategory{},
		&Cart{},
		&CartItem{},
	)

	// Adding foreignKey
	DB.Model(&Payment{}).AddForeignKey("currency_id", "currencies(id)", "RESTRICT", "RESTRICT")
	DB.Model(&Membership{}).AddForeignKey("user_id", "users(id)", "RESTRICT", "RESTRICT")
	DB.Model(&Membership{}).AddForeignKey("plan_id", "plans(id)", "RESTRICT", "RESTRICT")
	DB.Model(&Membership{}).AddForeignKey("payment_id", "payments(id)", "RESTRICT", "RESTRICT")
	DB.Model(&Product{}).AddForeignKey("currency_id", "currencies(id)", "RESTRICT", "RESTRICT")
	DB.Model(&Product{}).AddForeignKey("status_id", "statuses(id)", "RESTRICT", "RESTRICT")
	DB.Model(&Product{}).AddForeignKey("product_type_id", "product_types(id)", "RESTRICT", "RESTRICT")
	DB.Model(&ProductTag{}).AddForeignKey("tag_id", "tags(id)", "RESTRICT", "RESTRICT")
	DB.Model(&ProductTag{}).AddForeignKey("product_id", "products(id)", "RESTRICT", "RESTRICT")
	DB.Model(&ProductCategory{}).AddForeignKey("category_id", "categories(id)", "RESTRICT", "RESTRICT")
	DB.Model(&ProductCategory{}).AddForeignKey("product_id", "products(id)", "RESTRICT", "RESTRICT")
	DB.Model(&Cart{}).AddForeignKey("user_id", "users(id)", "RESTRICT", "RESTRICT")
	DB.Model(&CartItem{}).AddForeignKey("cart_id", "carts(id)", "RESTRICT", "RESTRICT")
	DB.Model(&CartItem{}).AddForeignKey("product_id", "products(id)", "RESTRICT", "RESTRICT")
}
