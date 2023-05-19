package withoutcard

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

type ratePairFullOKX struct {
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

func TriangleArbitrageOKX(asset string, amount float64) {
	start := time.Now()
	pair := GetPairFromJSONOKX(asset)

	for firstPairName, secondThirdPairs := range pair {
		GetTriangleArbitrageOKX(start, asset, firstPairName, secondThirdPairs, amount)
	}
}

func GetTriangleArbitrageOKX(t time.Time, asset, firstPairName string, pairsArray []string, amount float64) {
	for count, pp := range pairsArray {
		CalculateTriangleArbitrageOKX(t, count, asset, firstPairName, pp, amount)

	}
}

func CalculateTriangleArbitrageOKX(t time.Time, count int, asset, firstPairName, pp string, amount float64) {
	pairRate := strings.Split(pp, "|")
	pairRate = append(pairRate, firstPairName)
	pairRateValue := GetRatePairTriangleOKX(pairRate) //rate for pairs from OKX

	if pairRateValue[firstPairName][0] == 0.0 || pairRateValue[pairRate[0]][0] == 0 || pairRateValue[pairRate[1]][0] == 0 ||
		pairRateValue[firstPairName][2] == 0.0 || pairRateValue[pairRate[0]][2] == 0 || pairRateValue[pairRate[1]][2] == 0 {
		return
	}

	var transAmountFirst, transAmountSecond, result float64
	var rateFirst, rateSecond, rateThird float64

	// first step
	if strings.HasPrefix(firstPairName, asset) {
		transAmountFirst = amount * pairRateValue[firstPairName][0]
		rateFirst = pairRateValue[firstPairName][0]
	} else {
		transAmountFirst = amount / pairRateValue[firstPairName][2]
		rateFirst = pairRateValue[firstPairName][2]
	}
	tmpAssetSecond := strings.TrimSuffix(firstPairName, asset)

	//second steps
	if strings.HasPrefix(pairRate[0], tmpAssetSecond) {
		transAmountSecond = transAmountFirst * pairRateValue[pairRate[0]][0]
		rateSecond = pairRateValue[pairRate[0]][0]
	} else {
		transAmountSecond = transAmountFirst / pairRateValue[pairRate[0]][2]
		rateSecond = pairRateValue[pairRate[0]][2]
	}

	if strings.HasPrefix(pairRate[1], asset) {
		result = transAmountSecond / pairRateValue[pairRate[1]][2]
		rateThird = pairRateValue[pairRate[1]][2]
	} else {
		result = transAmountSecond * pairRateValue[pairRate[1]][0]
		rateThird = pairRateValue[pairRate[1]][0]
	}

	timeSince := time.Since(t)
	if result > amount {
		log.Printf("%d Asset: %s First Pair %s %.2f %.2f Secong Pair %s %.2f %.2f Third Pair %s %.2f %.2f - Result: %v",
			count, asset, firstPairName, rateFirst, transAmountFirst, pairRate[0], rateSecond,
			transAmountSecond, pairRate[1], rateThird, result, timeSince)
	}
}

func GetPairFromJSONOKX(asset string) map[string][]string {
	pair := ""
	switch asset {
	case "USDT":
		pair = fmt.Sprintf("data/dataokx/%s/%s_pairO_pairO_pairO.json", asset, asset)
	default:
		log.Printf("For %v don't have para\n", asset)
	}

	jsonfile, err := os.ReadFile(pair)
	if err != nil {
		panic(err)
	}
	var result map[string][]string
	_ = json.Unmarshal(jsonfile, &result)

	return result
}

func GetRatePairTriangleOKX(pair []string) map[string][4]float64 {
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

		res, err := sendRequestRatePairTriangleOKX(pair)

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

func sendRequestRatePairTriangleOKX(pair []string) (map[string][4]float64, error) {
	rate_pair := make(map[string][4]float64)
	for _, p := range pair {
		url := fmt.Sprintf("https://www.okx.com/priapi/v5/market/mult-tickers?instIds=%s", p)

		resp, err := http.Get(url)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Println("Error  in Body", err)
		}

		rj := ratePairFullOKX{}

		if err := json.Unmarshal(body, &rj); err != nil {
			log.Println("failed to parse JSON: %v", err)
			return nil, err
		}

		if len(rj.Data) > 0 {
			bid, _ := strconv.ParseFloat(rj.Data[0].BidPx, 64)
			bidVolume, _ := strconv.ParseFloat(rj.Data[0].BidSz, 64)
			ask, _ := strconv.ParseFloat(rj.Data[0].AskPx, 64)
			askVolume, _ := strconv.ParseFloat(rj.Data[0].AskPx, 64)

			rate_pair[p] = [4]float64{bid, bidVolume, ask, askVolume}
		}
	}

	return rate_pair, nil
}
