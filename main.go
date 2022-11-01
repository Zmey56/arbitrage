package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
)

type PostData struct {
	asset     string
	fiat      string
	page      int
	rows      int
	tradeType string
}

func main() {
	httpposturl := "https://p2p.binance.com/bapi/c2c/v2/friendly/c2c/adv/search"
	fmt.Println("HTTP JSON POST URL:", httpposturl)

	var jsonData = []byte(`{
		  "asset": "USDT",
		  "fiat": "RUB",
		  "page": 1,
		  "rows": 10,
		  "tradeType": "Buy"
	}`)
	request, error := http.NewRequest("POST", httpposturl, bytes.NewBuffer(jsonData))
	request.Header.Set("Content-Type", "application/json; charset=UTF-8")

	client := &http.Client{}
	response, error := client.Do(request)
	if error != nil {
		panic(error)
	}
	defer response.Body.Close()

	fmt.Println("response Status:", response.Status)
	fmt.Println("response Headers:", response.Header)
	body, _ := io.ReadAll(response.Body)
	fmt.Println("response Body:", string(body))

}
