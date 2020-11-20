package razorpay

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"github.com/spf13/viper"
)

// VerifySignature verifies razorpay payment signature
func VerifySignature(orderID, paymentID, signature string) bool {
	h := hmac.New(sha256.New, []byte(viper.GetString("razorpay_secret")))
	_, err := h.Write([]byte(fmt.Sprint(orderID, "|", paymentID)))
	if err != nil {
		return false
	}

	generated := hex.EncodeToString(h.Sum(nil))

	return generated == signature
}
