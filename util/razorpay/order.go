package razorpay

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

func CreateOrder(body map[string]interface{}) (map[string]interface{}, error) {
	buf := new(bytes.Buffer)

	err := json.NewEncoder(buf).Encode(&body)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", URL+"/orders", buf)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	AddAuthHeader(req)

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	resbody, _ := ioutil.ReadAll(resp.Body)

	order := map[string]interface{}{}

	err = json.Unmarshal(resbody, &order)

	if err != nil {
		return nil, err
	}

	return order, nil
}
