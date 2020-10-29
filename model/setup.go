package model

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

// DB - gorm DB
var DB *gorm.DB

// SetupDB is database setup
func SetupDB(DSN interface{}) {
	fmt.Println("connecting to database ...")

	var err error
	DB, err = gorm.Open(postgres.Open(viper.GetString("postgres.dsn")), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   "dp_",
			SingularTable: true,
		},
	})

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("connected to database ...")
}
