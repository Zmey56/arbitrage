package getinfobinance

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
)

type RatePair struct {
	LastUpdateId int64      `json:"lastUpdateId"`
	Bids         [][]string `json:"bids"`
	Asks         [][]string `json:"asks"`
}

func GetRatePairTriangleB(pair []string) map[string][4]float64 {
	result := make(map[string][4]float64)
	wg := sync.WaitGroup{}
	mu := sync.Mutex{}

	for _, p := range pair {
		wg.Add(1)
		go func(p string) {
			defer wg.Done()

			res, err := sendRequestRatePairTriangleB(p)
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

func sendRequestRatePairTriangleB(p string) ([4]float64, error) {
	ratePair := [4]float64{}

	url := fmt.Sprintf("https://www.binance.com/api/v3/depth?symbol=%s&limit=1", strings.ToUpper(p))

	resp, err := http.Get(url)
	if err != nil {
		log.Println(err)
		return ratePair, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error  in Body", err)
	}

	rj := RatePair{}

	if err := json.Unmarshal(body, &rj); err != nil {
		log.Printf("failed to parse JSON: %v, pair: %s", err, p)
		return ratePair, err
	}

	if len(rj.Bids) > 0 && len(rj.Asks) > 0 {
		bid, _ := strconv.ParseFloat(rj.Bids[0][0], 64)
		bidVolume, _ := strconv.ParseFloat(rj.Bids[0][1], 64)
		ask, _ := strconv.ParseFloat(rj.Asks[0][0], 64)
		askVolume, _ := strconv.ParseFloat(rj.Asks[0][1], 64)

		ratePair = [4]float64{bid, bidVolume, ask, askVolume}
	}

	return ratePair, nil
}
