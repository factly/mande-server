package config

import (
	"flag"
	"log"
)

// DSN dsn
var DSN string

// MeiliURL meili search server url
var MeiliURL string

// MeiliKey meili search key
var MeiliKey string

// RazorpayKey razorpay api key
var RazorpayKey string

// RazorpaySecret razorpay api secret
var RazorpaySecret string

// SetupVars setups all the config variables to run application
func SetupVars() {
	var dsn string
	var meili string
	var meiliKey string
	var razorpayKey string
	var razorpaySecret string

	flag.StringVar(&dsn, "dsn", "", "Database connection string")
	flag.StringVar(&meili, "meili", "", "Meili connection string")
	flag.StringVar(&meiliKey, "meiliKey", "", "Meili API Key string")
	flag.StringVar(&razorpayKey, "razorpayKey", "", "Razorpay API Key string")
	flag.StringVar(&razorpaySecret, "razorpaySecret", "", "Razorpay API Secret string")

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

	if razorpayKey == "" {
		log.Fatal("Please pass razorpayKey flag")
	}

	if razorpaySecret == "" {
		log.Fatal("Please pass razorpaySecret flag")
	}

	DSN = dsn
	MeiliURL = meili
	MeiliKey = meiliKey
	RazorpayKey = razorpayKey
	RazorpaySecret = razorpaySecret
}
