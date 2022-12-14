// Function for receiving data from the exchange Binance

package getinfobinance

import (
	"bytes"
	"encoding/json"
	"github.com/Zmey56/arbitrage/interact"
	"log"
	"net/http"
	"strconv"
)

func GetDataP2P(asset, fiat, tradeType string, paramUser interact.Parameters) AdvertiserAdv {
	count := 1
	transAmountFloat, err := strconv.ParseFloat(paramUser.TransAmount, 64)
	if err != nil {
		log.Println("Can't convert transAmount", err)
	}
	var jsonData = []byte(`{
				"payTypes": [],
				"rows": 10
			}`)
	var m map[string]interface{}
	err = json.Unmarshal(jsonData, &m)
	if err != nil {
		log.Println("Error", err)
	}
	m["asset"] = asset
	m["fiat"] = fiat
	m["tradeType"] = tradeType
	m["payTypes"] = paramUser.PayTypes
	m["page"] = count
	if paramUser.TransAmount != "" {
		m["transAmount"] = paramUser.TransAmount
	}
	if paramUser.PublisherType == "merchant" {
		m["publisherType"] = paramUser.PublisherType
	}

	newData, err := json.Marshal(m)
	if err != nil {
		log.Println(err)
	}
	newJsonData := []byte(string(newData))

	httpposturl := "https://p2p.binance.com/bapi/c2c/v2/friendly/c2c/adv/search"
	//fmt.Println("HTTP JSON POST URL:", httpposturl)

	request, error := http.NewRequest("POST", httpposturl, bytes.NewBuffer(newJsonData))
	request.Header.Set("Content-Type", "application/json; charset=UTF-8")

	client := &http.Client{}
	response, error := client.Do(request)
	if error != nil {
		panic(error)
	}
	defer response.Body.Close()

	resultadvertiseradv := ParsingJson(response.Body, tradeType, transAmountFloat)

	if resultadvertiseradv.Adv.Price != 0 {
		return resultadvertiseradv
	} else {
		var emptystruct AdvertiserAdv //return if don't find
		return emptystruct
	}
}
