// Function for receiving data from the exchange Binance

package getdata

import (
	"bytes"
	"encoding/json"
	"github.com/Zmey56/arbitrage/pkg/getinfobinance"
	"github.com/Zmey56/arbitrage/pkg/workingbinance"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"
)

func GetDataP2PBinance(asset, fiat, tradeType string, paramUser workingbinance.ParametersBinance) getinfobinance.AdvertiserAdv {
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

	resultadvertiseradv := getinfobinance.AdvertiserAdv{}

	for {
		defer func() {
			if r := recover(); r != nil {
				if r == "connection reset by peer" {
					log.Println("An error occured 'connection reset by peer', reconecting...")
					time.Sleep(time.Second * 1)
				} else {
					// Handling other errors
					log.Println("An error occured:", r)
					time.Sleep(time.Second * 1)
				}
			}
		}()

		resultadvertiseradv, err = requestOrdersP2P(newJsonData, tradeType, transAmountFloat)

		if err != nil {
			if err.Error() == "connection reset by peer" {
				// reconecting
				panic("connection reset by peer")
			} else {
				log.Println("Error:", err)
			}
		} else {
			break
		}
	}

	if len(resultadvertiseradv.Data) > 0 {
		return resultadvertiseradv
	} else {
		return getinfobinance.AdvertiserAdv{}
	}
}

func requestOrdersP2P(j []byte, tt string, taf float64) (getinfobinance.AdvertiserAdv, error) {
	httpposturl := "https://p2p.binance.com/bapi/c2c/v2/friendly/c2c/adv/search"
	//fmt.Println("HTTP JSON POST URL:", "https://p2p.binance.com/en/trade/"+string(j))
	request, err := http.NewRequest("POST", httpposturl, bytes.NewBuffer(j))
	if err != nil {
		return getinfobinance.AdvertiserAdv{}, err
	}
	request.Header.Set("Content-Type", "application/json; charset=UTF-8")

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return getinfobinance.AdvertiserAdv{}, err
	}

	defer response.Body.Close()

	return ParsingJson(response.Body, tt, taf), nil
}

func ParsingJson(r io.Reader, tradeType string, transAmount float64) getinfobinance.AdvertiserAdv {
	var result getinfobinance.AdvertiserAdv

	body, _ := io.ReadAll(r)
	err := json.Unmarshal([]byte(body), &result)

	if err != nil {
		log.Println("Error unmarshal json:", err)
	}

	return result
}
