package getinfohuobi

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

type ratePairFull struct {
	Ch     string `json:"ch"`
	Status string `json:"status"`
	Ts     int64  `json:"ts"`
	Tick   struct {
		Id      int64     `json:"id"`
		Version int64     `json:"version"`
		Open    float64   `json:"open"`
		Close   float64   `json:"close"`
		Low     float64   `json:"low"`
		High    float64   `json:"high"`
		Amount  float64   `json:"amount"`
		Vol     float64   `json:"vol"`
		Count   int       `json:"count"`
		Bid     []float64 `json:"bid"`
		Ask     []float64 `json:"ask"`
	} `json:"tick"`
}

func GetRatePairTriangleHuobi(pair []string) map[string][4]float64 {
	for {
		defer func() {
			if r := recover(); r != nil {
				if r == "connection reset by peer" {
					log.Println("An error occured 'connection reset by peer', reconecting...")
					time.Sleep(time.Second * 10)
				} else {
					// Handling other errors
					log.Println("An error occured:", r)
					time.Sleep(time.Second * 10)
					//return
				}
			}
		}()

		res, err := sendRequestRatePairTriangle(pair)

		if err != nil {
			if err.Error() == "connection reset by peer" {
				// reconecting
				panic("connection reset by peer")
			} else {
				log.Println(err)
				panic("error")
			}
		} else {
			return res
		}

	}
}

func sendRequestRatePairTriangle(pair []string) (map[string][4]float64, error) {
	rate_pair := make(map[string][4]float64)
	for _, p := range pair {
		//log.Println("pair", p)

		url := fmt.Sprintf("https://api.huobi.pro/market/detail/merged?symbol=%s", strings.ToLower(p))

		resp, err := http.Get(url)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Println("Error  in Body", err)
		}

		rj := ratePairFull{}

		if err := json.Unmarshal(body, &rj); err != nil {
			log.Println("failed to parse JSON: %v", err)
			return nil, err
		}
		//log.Println(rj.Tick.Bid, " - ", rj.Tick.Ask)

		if len(rj.Tick.Bid) > 0 && len(rj.Tick.Ask) > 0 {
			bid := rj.Tick.Bid[0]
			bidVolume := rj.Tick.Bid[1]
			ask := rj.Tick.Ask[0]
			askVolume := rj.Tick.Ask[1]

			rate_pair[p] = [4]float64{bid, bidVolume, ask, askVolume}
			//log.Println(p, " - ", rate_pair[p], " - ", [4]float64{bid, bidVolume, ask, askVolume})
		}
	}
	//for l, m := range rate_pair {
	//	fmt.Println(l, " - ", m)
	//}

	return rate_pair, nil
}
