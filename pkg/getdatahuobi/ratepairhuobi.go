//get rate from binance and return mean rate from json

package getdatahuobi

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"sync"
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

type HuobiRateVer2 struct {
	Tick struct {
		Close float64 `json:"close"`
	} `json:"tick"`
}

func GetRatePairHuobiVer2(pair []string) map[string]float64 {
	rate_pair := make(map[string]float64)
	var wg sync.WaitGroup

	for _, p := range pair {
		wg.Add(1)
		go func(pair string) {
			defer wg.Done()

			res, err := SendRequesrRatePairVer2([]string{pair})

			if err != nil {
				log.Println("Error for pair", pair, ":", err)
			} else {
				rate_pair[pair] = res[pair]
			}
		}(p)
	}

	wg.Wait()
	return rate_pair
}

func SendRequesrRatePairVer2(pair []string) (map[string]float64, error) {
	rate_pair := make(map[string]float64)
	for _, p := range pair {
		url := fmt.Sprintf("https://api.huobi.pro/market/detail/merged?symbol=%s", strings.ToLower(p))

		resp, err := http.Get(url)
		if err != nil {
			return rate_pair, err
		}
		defer resp.Body.Close()
		body, _ := io.ReadAll(resp.Body)

		ratejson := HuobiRateVer2{}

		if err := json.Unmarshal(body, &ratejson); err != nil {
			return rate_pair, err
		}

		rate_pair[p] = ratejson.Tick.Close
	}
	return rate_pair, nil
}
