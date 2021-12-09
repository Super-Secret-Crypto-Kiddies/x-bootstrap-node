package p2p

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type APIResponse struct {
	Query string `json:"query"`
}

func GetIPAddress() string {
	resp, err := http.Get("http://ip-api.com/json/")

	if err != nil {
		panic(err)
	}

	body, _ := ioutil.ReadAll(resp.Body)

	var apiResp APIResponse
	json.Unmarshal(body, &apiResp)

	return apiResp.Query
}
