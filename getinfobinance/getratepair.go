//get rate from binance and return mean rate from json

package getinfobinance

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"
)

type ratejson struct {
	LastUpdateID int64      `json:"lastUpdateId"`
	Bids         [][]string `json:"bids"`
	Asks         [][]string `json:"asks"`
}

func GetRatePair(pair []string) map[string]float64 {
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
		url := fmt.Sprintf("https://www.binance.com/api/v3/depth?symbol=%s&limit=1", p)

		resp, err := http.Get(url)
		if err != nil {
			return rate_pair, err
		}
		defer resp.Body.Close()
		body, _ := io.ReadAll(resp.Body)
		//log.Println(p, " - ", string(body))

		rj := ratejson{}

		if err := json.Unmarshal(body, &rj); err != nil {
			return rate_pair, err
		}

		//log.Println(len(rj.Bids), len(rj.Asks))

		if len(rj.Bids) > 0 && len(rj.Asks) > 0 {
			bids, _ := strconv.ParseFloat(rj.Bids[0][0], 64)
			asks, _ := strconv.ParseFloat(rj.Asks[0][0], 64)

			rate_pair[p] = (bids + asks) / 2.0
		}
	}
	return rate_pair, nil
}
