package getdataokx

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

type AllRatePairOKX struct {
	Code string `json:"code"`
	Msg  string `json:"msg"`
	Data []struct {
		InstType  string `json:"instType"`
		InstId    string `json:"instId"`
		Last      string `json:"last"`
		LastSz    string `json:"lastSz"`
		AskPx     string `json:"askPx"`
		AskSz     string `json:"askSz"`
		BidPx     string `json:"bidPx"`
		BidSz     string `json:"bidSz"`
		Open24H   string `json:"open24h"`
		High24H   string `json:"high24h"`
		Low24H    string `json:"low24h"`
		VolCcy24H string `json:"volCcy24h"`
		Vol24H    string `json:"vol24h"`
		Ts        string `json:"ts"`
		SodUtc0   string `json:"sodUtc0"`
		SodUtc8   string `json:"sodUtc8"`
	} `json:"data"`
}

func GetRatePairTriangleO(pair []string, fiat string) map[string][4]float64 {
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

		res, err := sendRequesrRatePairOKXWithout(pair, fiat)

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

func sendRequesrRatePairOKXWithout(pair []string, fiat string) (map[string][4]float64, error) {
	rate_pair := make(map[string][4]float64)

	count := 0

	for count < 10 {

		var pairUpdate []string

		if count != 0 {
			for _, p := range pair {
				_, ok := rate_pair[p]
				if !ok {
					pairUpdate = append(pairUpdate, p)
				}
			}
		} else {
			pairUpdate = pair
		}

		if len(pairUpdate) < 1 {
			break
		}

		var pairDef []string

		for _, p := range pairUpdate {
			pairDef = append(pairDef, addDefis(p, fiat))
		}

		url := fmt.Sprintf("https://www.okx.com/priapi/v5/market/mult-tickers?instIds=%s", strings.Join(pairDef, ","))

		resp, err := makeRequest(url)

		// Make sure the response body is closed when we're done with it
		defer resp.Body.Close()

		body, _ := io.ReadAll(resp.Body)

		ratejson := AllRatePairOKX{}

		if err := json.Unmarshal(body, &ratejson); err != nil {
			return nil, errors.New("Can't unmarshal from rate pair")
		}
		if err != nil {
			log.Println("Error:", err)
			return nil, errors.New("Can't convert rate pair on OKX")
		}

		for _, value := range ratejson.Data {
			var tmp [4]float64

			bid, _ := strconv.ParseFloat(value.BidPx, 64)
			bidVolume, _ := strconv.ParseFloat(value.BidSz, 64)
			ask, _ := strconv.ParseFloat(value.AskPx, 64)
			askVolume, _ := strconv.ParseFloat(value.AskSz, 64)

			tmp = [4]float64{bid, bidVolume, ask, askVolume}

			pairName := strings.Join(strings.Split(value.InstId, "-"), "")

			rate_pair[strings.ToLower(pairName)] = tmp
		}

		count++

	}

	return rate_pair, nil

}

func addDefis(pair, fiat string) string {
	if strings.HasPrefix(pair, fiat) {
		secondCoin := pair[len(fiat):]
		return fmt.Sprintf("%s-%s", strings.ToUpper(fiat), strings.ToUpper(secondCoin))

	} else {
		firstCoin := pair[:len(pair)-len(fiat)]
		return fmt.Sprintf("%s-%s", strings.ToUpper(firstCoin), strings.ToUpper(fiat))
	}
}

func GetPairOKXWithout(fiat string) []string {
	pair := ""
	switch fiat {
	case "USDT":
		pair = fmt.Sprintf("data/dataokx/%s/%s_pair_without.json", fiat, fiat)
	default:
		log.Printf("For %v don't have para\n", fiat)
	}

	jsonfile, err := os.ReadFile(pair)
	if err != nil {
		panic(err)
	}
	var result []string
	_ = json.Unmarshal(jsonfile, &result)

	return result
}
