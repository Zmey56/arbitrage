// Function for receiving data from the exchange Binance

package getdata

import (
	"bytes"
	"encoding/json"
	"github.com/Zmey56/arbitrage/pkg/getinfobinance"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

// GetDataP2PBinanceVer2 get all data for first page from Binance
func GetDataP2PBinanceVer2(asset, fiat, tradeType string, page int) getinfobinance.Binance {
	var jsonData = []byte(`{
                "payTypes": [],
                "rows": 10,
                "publisherType": null
            }`)
	var m map[string]interface{}
	err := json.Unmarshal(jsonData, &m)
	if err != nil {
		log.Println("Error with Unmarshal", err)
	}
	m["asset"] = asset
	m["fiat"] = fiat
	m["tradeType"] = tradeType
	m["page"] = page

	newData, err := json.Marshal(m)
	if err != nil {
		log.Println("Error with Marshal", err)
	}
	newJsonData := []byte(string(newData))

	resultadvertiseradv := getinfobinance.Binance{}

	for retries := 0; retries < 3; retries++ {
		resultadvertiseradv, err = requestOrdersP2PVer2(newJsonData, tradeType)
		if err != nil {
			if strings.Contains(err.Error(), "connection reset by peer") {
				log.Printf("Error: %s, retrying in 1 second...\n", err)
				time.Sleep(time.Second * 1)
				continue
			} else {
				log.Println("Error:", err)
				break
			}
		} else {
			break
		}
	}

	if len(resultadvertiseradv.Data) > 0 {
		return resultadvertiseradv
	} else {
		return getinfobinance.Binance{}
	}
}

func requestOrdersP2PVer2(j []byte, tt string) (getinfobinance.Binance, error) {
	httpposturl := "https://p2p.binance.com/bapi/c2c/v2/friendly/c2c/adv/search"
	request, err := http.NewRequest("POST", httpposturl, bytes.NewBuffer(j))
	if err != nil {
		return getinfobinance.Binance{}, err
	}
	request.Header.Set("Content-Type", "application/json; charset=UTF-8")
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return getinfobinance.Binance{}, err
	}

	defer response.Body.Close()

	dateStr := response.Header.Get("Date")
	if dateStr == "" {
		log.Println("Data header not found")
		return getinfobinance.Binance{}, err
	}
	layout := http.TimeFormat
	date, err := time.Parse(layout, dateStr)
	if err != nil {
		log.Println("Invalid Date header")
		return getinfobinance.Binance{}, err
	}
	unixTime := date.Unix()

	return parsingJsonVer2(response.Body, tt, unixTime), nil
}

func parsingJsonVer2(r io.Reader, tradeType string, ut int64) getinfobinance.Binance {
	var result getinfobinance.Binance

	body, _ := io.ReadAll(r)
	err := json.Unmarshal([]byte(body), &result)

	if err != nil {
		log.Println("Error unmarshal json from URL:", err, "/n")
	}

	result.TimeData = ut

	return result
}
