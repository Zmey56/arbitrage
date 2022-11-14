// Function for receiving data from the exchange Binance

package p2pbinance

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func GetDataP2P(asset, fiat, tradeType string, payTypes []string) AdvertiserAdv {
	var jsonData = []byte(`{
			"page": 1,
			"payTypes": [],
		  	"rows": 10
	}`)
	var m map[string]interface{}
	err := json.Unmarshal(jsonData, &m)
	if err != nil {
		log.Println("Error", err)
	}
	m["asset"] = asset
	m["fiat"] = fiat
	m["tradeType"] = tradeType
	m["payTypes"] = payTypes

	newData, err := json.Marshal(m)
	if err != nil {
		fmt.Println(err)
	}
	newJsonData := []byte(string(newData))

	httpposturl := "https://p2p.binance.com/bapi/c2c/v2/friendly/c2c/adv/search"
	fmt.Println("HTTP JSON POST URL:", httpposturl)

	request, error := http.NewRequest("POST", httpposturl, bytes.NewBuffer(newJsonData))
	request.Header.Set("Content-Type", "application/json; charset=UTF-8")

	client := &http.Client{}
	response, error := client.Do(request)
	if error != nil {
		panic(error)
	}
	defer response.Body.Close()

	tmp, _ := ParsingJson(response.Body)

	return (tmp)
}
