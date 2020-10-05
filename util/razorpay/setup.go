package razorpay

import (
	"encoding/base64"
	"fmt"
	"net/http"

	"github.com/factly/data-portal-server/config"
)

var URL string = "https://api.razorpay.com/v1"

func AddAuthHeader(request *http.Request) {
	authHeader := fmt.Sprint(config.RazorpayKey, ":", config.RazorpaySecret)
	authHeaderEnc := base64.StdEncoding.EncodeToString([]byte(authHeader))

	header := fmt.Sprint("Basic ", authHeaderEnc)
	request.Header.Add("Authorization", header)
}
