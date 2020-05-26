package model

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

// SetupDB is database setup
func SetupDB() {
	env := os.Getenv("ENVIRONMENT")
	if "" == env {
		env = "local"
	}
	envFileName := ".env." + env
	err := godotenv.Load(envFileName)
	if err != nil {
		log.Fatal("error loading .env.local file")
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

	DB.SingularTable(true)

	// adding default prefix to table name
	gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
		return "dp_" + defaultTableName
	}

	// db migrations
	DB.AutoMigrate(
		&Currency{},
		&Payment{},
		&Membership{},
		&Plan{},
		&User{},
		&Product{},
		&ProductType{},
		&Tag{},
		&ProductTag{},
		&Category{},
		&ProductCategory{},
		&Cart{},
		&CartItem{},
		&Order{},
		&OrderItem{},
	)

	// Adding foreignKey
	DB.Model(&Payment{}).AddForeignKey("currency_id", "dp_currency(id)", "RESTRICT", "RESTRICT")
	DB.Model(&Membership{}).AddForeignKey("user_id", "dp_user(id)", "RESTRICT", "RESTRICT")
	DB.Model(&Membership{}).AddForeignKey("plan_id", "dp_plan(id)", "RESTRICT", "RESTRICT")
	DB.Model(&Membership{}).AddForeignKey("payment_id", "dp_payment(id)", "RESTRICT", "RESTRICT")
	DB.Model(&Product{}).AddForeignKey("currency_id", "dp_currency(id)", "RESTRICT", "RESTRICT")
	DB.Model(&Product{}).AddForeignKey("product_type_id", "dp_product_type(id)", "RESTRICT", "RESTRICT")
	DB.Model(&ProductTag{}).AddForeignKey("tag_id", "dp_tag(id)", "RESTRICT", "RESTRICT")
	DB.Model(&ProductTag{}).AddForeignKey("product_id", "dp_product(id)", "RESTRICT", "RESTRICT")
	DB.Model(&ProductCategory{}).AddForeignKey("category_id", "dp_category(id)", "RESTRICT", "RESTRICT")
	DB.Model(&ProductCategory{}).AddForeignKey("product_id", "dp_product(id)", "RESTRICT", "RESTRICT")
	DB.Model(&Cart{}).AddForeignKey("user_id", "dp_user(id)", "RESTRICT", "RESTRICT")
	DB.Model(&CartItem{}).AddForeignKey("cart_id", "dp_cart(id)", "RESTRICT", "RESTRICT")
	DB.Model(&CartItem{}).AddForeignKey("product_id", "dp_product(id)", "RESTRICT", "RESTRICT")
	DB.Model(&Order{}).AddForeignKey("payment_id", "dp_payment(id)", "RESTRICT", "RESTRICT")
	DB.Model(&Order{}).AddForeignKey("cart_id", "dp_cart(id)", "RESTRICT", "RESTRICT")
	DB.Model(&OrderItem{}).AddForeignKey("product_id", "dp_product(id)", "RESTRICT", "RESTRICT")
	DB.Model(&OrderItem{}).AddForeignKey("order_id", "dp_order(id)", "RESTRICT", "RESTRICT")
}
