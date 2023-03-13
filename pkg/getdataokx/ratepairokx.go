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
		//log.Println("Result into function", res)

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
	//log.Println("Pair into function", pair)
	for _, p := range pair {
		//log.Println("PAIR", p)
		url := fmt.Sprintf("https://www.okx.com/priapi/v5/market/mult-tickers?instIds=%s", p)
		//log.Printf("URL for pair %s is %s", p, url)
		resp, err := makeRequest(url)

		//req, err := http.NewRequest("GET", url, nil)
		//if err != nil {
		//	fmt.Println("Error creating HTTP request:", err)
		//	return nil, err
		//}
		//req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.36")
		//req.Header.Set("Accept-Language", "en-US,en;q=0.5")
		//
		//// Send the HTTP request using the default HTTP client
		//resp, err := http.DefaultClient.Do(req)
		//if err != nil {
		//	fmt.Println("Error sending HTTP request:", err)
		//	return nil, err
		//}
		//log.Println("StatusCode", resp.StatusCode)
		//log.Println("RESP BODY", resp.Body)

		// Make sure the response body is closed when we're done with it
		defer resp.Body.Close()

		//resp, err := http.Get(url)
		//if err != nil {
		//	return rate_pair, err
		//}
		//defer resp.Body.Close()
		body, _ := io.ReadAll(resp.Body)

		ratejson := RatePairOKX{}

		if err := json.Unmarshal(body, &ratejson); err != nil {
			return nil, errors.New("Can't unmarshal from rate pair")
		}
		//log.Println("ratejson", ratejson, "PAIR", p, "URL", url)

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
		// Задержка перед повтором запроса
		retryAfter := 5 * time.Nanosecond
		//log.Println("resp.Header", resp.Header)
		//retryAfter := resp.Header.Get("Retry-After")
		//log.Println("retryAfter", retryAfter)
		//retryAfterDuration, err := time.ParseDuration(retryAfter + "ns")
		if err != nil {
			return nil, err
		}
		time.Sleep(retryAfter)

		// Повтор запроса
		return makeRequest(url)
	}

	return resp, nil
}

// Get MAP crypto and pair
func GetPairFromJSONOKX(fiat string) map[string][]string {
	file_path := fmt.Sprintf("data/dataokx/%s/%s_pair.json", fiat, fiat)
	log.Println(file_path)
	file, _ := os.Open(file_path)
	defer file.Close()
	decoder := json.NewDecoder(file)
	var data map[string][]string
	decoder.Decode(&data)
	return (data)
}
