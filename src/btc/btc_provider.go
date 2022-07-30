package btc

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"os"
	"strconv"
)

type BtcResponce struct {
	Price string `json:"price"`
	Mins  int    `json:"mins"`
}

func GetRate() int {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://api.binance.com/api/v3/avgPrice?symbol=BTCUAH", nil)
	if err != nil {
		log.Print(err)
		os.Exit(1)
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Print(err)
		os.Exit(1)
	}
	var responce BtcResponce
	respBody, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal([]byte(respBody), &responce)
	btcRate, err := strconv.ParseFloat(responce.Price, 64)
	if err != nil {
		log.Print(err)
	}
	return int(math.Round(btcRate))
}
