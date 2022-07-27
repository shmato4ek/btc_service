package main

import (
	"btc_service/src/persistance"
	"fmt"
)

func main() {
	database := persistance.New("D:/BTC service/src")
	fmt.Println(database.Exists("m.shmatko8@gmail.com"))
}
