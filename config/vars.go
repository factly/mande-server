package config

import (
	"log"

	"github.com/spf13/viper"
)

// SetupVars setups all the config variables to run application
func SetupVars() {
	var configPath string

	viper.AddConfigPath(".")
	viper.SetConfigName("config")
	viper.SetEnvPrefix("mande_")
	viper.AutomaticEnv()

	viper.SetConfigFile(configPath)

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal("config file not found...")
	}
	if !viper.IsSet("database_host") {
		log.Fatal("please provide database_host config param")
	}

	if !viper.IsSet("database_user") {
		log.Fatal("please provide database_user config param")
	}

	if !viper.IsSet("database_name") {
		log.Fatal("please provide database_name config param")
	}

	if !viper.IsSet("database_password") {
		log.Fatal("please provide database_password config param")
	}

	if !viper.IsSet("database_port") {
		log.Fatal("please provide database_port config param")
	}

	if !viper.IsSet("database_ssl_mode") {
		log.Fatal("please provide database_ssl_mode config param")
	}

	if !viper.IsSet("meili_url") {
		log.Fatal("please provide meili_url in config file")
	}

	if !viper.IsSet("meili_key") {
		log.Fatal("please provide meili_key in config file")
	}

	if !viper.IsSet("razorpay_key") {
		log.Fatal("please provide razorpay_key in config file")
	}

	if !viper.IsSet("razorpay_secret") {
		log.Fatal("please provide razorpay_secret in config file")
	}
}
