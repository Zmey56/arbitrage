//get rate from binance and return mean rate from json

package getinfobinance

import (
	"encoding/json"
	"fmt"
	"io"
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
	rate_pair := make(map[string]float64)
	for _, p := range pair {
		url := fmt.Sprintf("https://www.binance.com/api/v3/depth?symbol=%s&limit=1", p)
		resp, err := http.Get(url)
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()
		body, _ := io.ReadAll(resp.Body)
		//log.Println(p, " - ", string(body))

		rj := ratejson{}

		if err := json.Unmarshal(body, &rj); err != nil {
			panic(err)
		}

		//log.Println(len(rj.Bids), len(rj.Asks))

		if len(rj.Bids) > 0 && len(rj.Asks) > 0 {
			bids, _ := strconv.ParseFloat(rj.Bids[0][0], 64)
			asks, _ := strconv.ParseFloat(rj.Asks[0][0], 64)

			rate_pair[p] = (bids + asks) / 2.0
		}
		time.Sleep(1 * time.Nanosecond)
	}
	return rate_pair
}
