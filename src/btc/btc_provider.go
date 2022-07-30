package btc

import (
	"encoding/json"
	"io/ioutil"
	"math"
	"net/http"
	"strconv"
)

type BtcResponce struct {
	Price string `json:"price"`
	Mins  int    `json:"mins"`
}

func GetRate(apiName string) (int, error) {
	client := &http.Client{}
	req, _ := http.NewRequest("GET", apiName, nil)

	resp, err := client.Do(req)
	var responce BtcResponce
	respBody, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal([]byte(respBody), &responce)
	btcRate, _ := strconv.ParseFloat(responce.Price, 64)
	return int(math.Round(btcRate)), err
}
