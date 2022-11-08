// Function for receiving data from the exchange Binance

package p2pbinance

import (
	"encoding/json"
	"fmt"
	"log"
)

func GetDataP2P(asset, fiat, tradeType string, payTypes []string) []byte {
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
	fmt.Println(jsonData)
	m["asset"] = asset
	m["fiat"] = fiat
	m["tradeType"] = tradeType

	newData, err := json.Marshal(m)
	if err != nil {
		fmt.Println(err)
	}
	newJsonData := []byte(string(newData))

	return (newJsonData)
}

// get information about payment methods and fiat
func GetDataP2PBank(fiat string) []byte {
	var jsonData = []byte(`{
	}`)
	var m map[string]interface{}
	err := json.Unmarshal(jsonData, &m)
	if err != nil {
		log.Println("Error", err)
	}
	m["fiat"] = fiat

	newData, err := json.Marshal(m)
	if err != nil {
		fmt.Println(err)
	}
	newJsonData := []byte(string(newData))

	return (newJsonData)
}

// get information about payment methods and fiat
func GetDataP2PCrypto(fiat string) []byte {
	var jsonData = []byte(`{
	}`)
	var m map[string]interface{}
	err := json.Unmarshal(jsonData, &m)
	if err != nil {
		log.Println("Error", err)
	}
	m["fiat"] = fiat

	newData, err := json.Marshal(m)
	if err != nil {
		fmt.Println(err)
	}
	newJsonData := []byte(string(newData))

	return (newJsonData)
}
