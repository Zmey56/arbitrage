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

type RatePair struct {
	LastUpdateId int64      `json:"lastUpdateId"`
	Bids         [][]string `json:"bids"`
	Asks         [][]string `json:"asks"`
}

func GetRatePairTriangle(pair []string) map[string][4]float64 {
	for {
		defer func() {
			if r := recover(); r != nil {
				if r == "connection reset by peer" {
					log.Println("An error occured 'connection reset by peer', reconecting...")
					time.Sleep(time.Second * 5)
				} else {
					// Handling other errors
					log.Println("An error occured:", r)
					time.Sleep(time.Second * 5)
					//return
				}
			}
		}()

		res, err := SendRequestRatePairTriangle(pair)

		if err != nil {
			if err.Error() == "connection reset by peer" {
				// reconecting
				panic("connection reset by peer")
			} else {
				panic("error")
			}
		} else {
			return res
		}

	}
}

func SendRequestRatePairTriangle(pair []string) (map[string][4]float64, error) {
	rate_pair := make(map[string][4]float64)
	for _, p := range pair {
		url := fmt.Sprintf("https://www.binance.com/api/v3/depth?symbol=%s&limit=1", p)

		resp, err := http.Get(url)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Println("Error  in Body", err)
		}

		rj := RatePair{}

		if err := json.Unmarshal(body, &rj); err != nil {
			return nil, err
		}

		if len(rj.Bids) > 0 && len(rj.Asks) > 0 {
			bid, _ := strconv.ParseFloat(rj.Bids[0][0], 64)
			bidVolume, _ := strconv.ParseFloat(rj.Bids[0][1], 64)
			ask, _ := strconv.ParseFloat(rj.Asks[0][0], 64)
			askVolume, _ := strconv.ParseFloat(rj.Asks[0][1], 64)

			rate_pair[p] = [4]float64{bid, bidVolume, ask, askVolume}
			//log.Println(p, " - ", rate_pair[p], " - ", [4]float64{bid, bidVolume, ask, askVolume})
		}
	}
	return rate_pair, nil
}
