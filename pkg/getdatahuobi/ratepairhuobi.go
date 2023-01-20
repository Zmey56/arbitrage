//get rate from binance and return mean rate from json

package getdatahuobi

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

func GetRatePairHuobi(pair []string) map[string]float64 {
	for {
		defer func() {
			if r := recover(); r != nil {
				if r == "connection reset by peer" {
					log.Println("An error occured 'connection reset by peer', reconecting...")
					time.Sleep(time.Second * 1)
				} else {
					// Handling other errors
					log.Println("An error occured:", r)
					return
				}
			}
		}()

		res, err := SendRequesrRatePair(pair)

		if err != nil {
			if err.Error() == "connection reset by peer" {
				// reconecting
				panic("connection reset by peer")
			} else {
				log.Println("Error:", err)
			}
		} else {
			return res
		}

	}
}

func SendRequesrRatePair(pair []string) (map[string]float64, error) {
	rate_pair := make(map[string]float64)
	for _, p := range pair {
		url := fmt.Sprintf("https://api.huobi.pro/market/detail/merged?symbol=%s", p)

		resp, err := http.Get(url)
		if err != nil {
			return rate_pair, err
		}
		defer resp.Body.Close()
		body, _ := io.ReadAll(resp.Body)

		ratejson := HuobiRate{}

		if err := json.Unmarshal(body, &ratejson); err != nil {
			return rate_pair, err
		}

		rate_pair[p] = ratejson.Tick.Close
	}
	return rate_pair, nil
}
