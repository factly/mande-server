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
	viper.SetEnvPrefix("mande")
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

	if !viper.IsSet("super_organisation_title") {
		log.Fatal("please provide super_organisation_title config param")
	}

	if !viper.IsSet("default_user_email") {
		log.Fatal("please provide default_user_email config param")
	}

	if !viper.IsSet("default_user_password") {
		log.Fatal("please provide default_user_password config param")
	}

	if !viper.IsSet("keto_url") {
		log.Fatal("please provide keto_url config param")
	}

	if !viper.IsSet("kavach_url") {
		log.Fatal("please provide kavach_url config param")
	}

	if !viper.IsSet("kratos_public_url") {
		log.Fatal("please provide kratos_public_url config param")
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
