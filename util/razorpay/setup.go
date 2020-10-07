package razorpay

import (
	"github.com/factly/data-portal-server/config"
	"github.com/razorpay/razorpay-go"
)

// Client client for razorpay
var Client *razorpay.Client

// SetupClient setups the client with key and secret
func SetupClient() {
	Client = razorpay.NewClient(config.RazorpayKey, config.RazorpaySecret)
}
