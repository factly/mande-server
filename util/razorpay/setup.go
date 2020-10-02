package razorpay

import (
	"github.com/factly/data-portal-server/config"
	razorpay "github.com/razorpay/razorpay-go"
)

// Client client for razorpay
var Client *razorpay.Client

// SetupRazorpay setups razorpay client to request razorpay api
func SetupRazorpay() {
	Client = razorpay.NewClient(config.RazorpayKey, config.RazorpaySecret)
}
