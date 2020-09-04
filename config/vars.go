package config

import (
	"flag"
	"log"
)

// DSN dsn
var DSN string

// MeiliURL meili search server url
var MeiliURL string

// MeiliKey meili search server url
var MeiliKey string

// SetupVars setups all the config variables to run application
func SetupVars() {
	var dsn string
	var meili string
	var meiliKey string

	flag.StringVar(&dsn, "dsn", "", "Database connection string")
	flag.StringVar(&meili, "meili", "", "Meili connection string")
	flag.StringVar(&meiliKey, "meiliKey", "", "Meili API Key string")

	flag.Parse()

	if dsn == "" {
		log.Fatal("Please pass dsn flag")
	}

	if meili == "" {
		log.Fatal("Please pass meili flag")
	}

	if meiliKey == "" {
		log.Fatal("Please pass meiliKey flag")
	}

	DSN = dsn
	MeiliURL = meili
	MeiliKey = meiliKey
}
