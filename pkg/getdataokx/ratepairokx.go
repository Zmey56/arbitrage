package getdataokx

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

func GetRatePairOKX(pair []string) map[string]float64 {
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

		res, err := sendRequesrRatePairOKX(pair)

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

func sendRequesrRatePairOKX(pair []string) (map[string]float64, error) {
	rate_pair := make(map[string]float64)
	for _, p := range pair {
		url := fmt.Sprintf("https://www.okx.com/priapi/v5/market/mult-tickers?instIds=%s", p)

		resp, err := makeRequest(url)

		// Make sure the response body is closed when we're done with it
		defer resp.Body.Close()

		body, _ := io.ReadAll(resp.Body)

		ratejson := RatePairOKX{}

		if err := json.Unmarshal(body, &ratejson); err != nil {
			return nil, errors.New("Can't unmarshal from rate pair")
		}

		rate_pair[p], err = strconv.ParseFloat(ratejson.Data[0].Last, 64)

		if err != nil {
			log.Println("Error:", err)
			return nil, errors.New("Can't convert rate pair on OKX")
		}
	}
	return rate_pair, nil
}

func makeRequest(url string) (*http.Response, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == http.StatusTooManyRequests {
		retryAfter := 5 * time.Nanosecond
		if err != nil {
			return nil, err
		}
		time.Sleep(retryAfter)

		return makeRequest(url)
	}

	return resp, nil
}

// Get MAP crypto and pair
func GetPairFromJSONOKX(fiat string) map[string][]string {
	file_path := fmt.Sprintf("data/dataokx/%s/%s_pair.json", fiat, fiat)
	file, _ := os.Open(file_path)
	defer file.Close()
	decoder := json.NewDecoder(file)
	var data map[string][]string
	decoder.Decode(&data)
	return (data)
}
