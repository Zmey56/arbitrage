package getdatahuobi

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"sync"
)

func GetRatePairTriangleH(pair []string) map[string][4]float64 {
	result := make(map[string][4]float64)
	wg := sync.WaitGroup{}
	mu := sync.Mutex{}

	for _, p := range pair {
		wg.Add(1)
		go func(p string) {
			defer wg.Done()

			res, err := sendRequestRatePairTriangleH(p)
			if err != nil {
				log.Println(err)
				return
			}

			mu.Lock()
			result[p] = res
			mu.Unlock()
		}(p)
	}

	wg.Wait()

	return result
}

func sendRequestRatePairTriangleH(p string) ([4]float64, error) {
	rate_pair := [4]float64{}
	url := fmt.Sprintf("https://api.huobi.pro/market/detail/merged?symbol=%s", strings.ToLower(p))

	resp, err := http.Get(url)
	if err != nil {
		return rate_pair, err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)

	rj := HuobiRate{}

	if err := json.Unmarshal(body, &rj); err != nil {
		log.Printf("failed to parse JSON: %v, pair %s", err, p)
		return rate_pair, err
	}

	if len(rj.Tick.Bid) > 0 && len(rj.Tick.Ask) > 0 {
		bid := rj.Tick.Bid[0]
		bidVolume := rj.Tick.Bid[1]
		ask := rj.Tick.Ask[0]
		askVolume := rj.Tick.Ask[1]

		rate_pair = [4]float64{bid, bidVolume, ask, askVolume}

	}
	return rate_pair, nil
}
