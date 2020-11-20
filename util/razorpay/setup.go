package razorpay

import (
	"github.com/razorpay/razorpay-go"
	"github.com/spf13/viper"
)

// Client client for razorpay
var Client *razorpay.Client

// SetupClient setups the client with key and secret
func SetupClient() {
	Client = razorpay.NewClient(viper.GetString("razorpay_key"), viper.GetString("razorpay_secret"))
}
