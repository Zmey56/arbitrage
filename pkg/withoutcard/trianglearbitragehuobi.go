package withoutcard

import (
	"encoding/json"
	"fmt"
	"github.com/Zmey56/arbitrage/pkg/getinfohuobi"
	"log"
	"os"
	"strings"
	"time"
)

func TriangleArbitrageHuobi(asset string, amount float64) {
	start := time.Now()
	pair := GetPairFromJSONHuobi(asset)

	for firstPairName, secondThirdPairs := range pair {
		GetTriangleArbitrageHuobi(start, asset, firstPairName, secondThirdPairs, amount)
	}
}

func GetTriangleArbitrageHuobi(t time.Time, asset, firstPairName string, pairsArray []string, amount float64) {
	for count, pp := range pairsArray {
		CalculateTriangleArbitrageHuobi(t, count, asset, firstPairName, pp, amount)

	}
}

func CalculateTriangleArbitrageHuobi(t time.Time, count int, asset, firstPairName, pp string, amount float64) {
	pairRate := strings.Split(pp, "|")
	pairRate = append(pairRate, firstPairName)
	pairRateValue := getinfohuobi.GetRatePairTriangleHuobi(pairRate) //rate for pairs from Huobi
	//pairRateValueB := getinfobinance.GetRatePairTriangle([]string{firstPairName}) //rate for pair from Binance

	if pairRateValue[firstPairName][0] == 0.0 || pairRateValue[pairRate[0]][0] == 0 || pairRateValue[pairRate[1]][0] == 0 ||
		pairRateValue[firstPairName][2] == 0.0 || pairRateValue[pairRate[0]][2] == 0 || pairRateValue[pairRate[1]][2] == 0 {
		return
	}

	firstPairNameLower := strings.ToLower(firstPairName)
	assetLower := strings.ToLower(asset)
	var transAmountFirst, transAmountSecond, result float64
	var rateFirst, rateSecond, rateThird float64

	// first step
	if strings.HasPrefix(firstPairNameLower, assetLower) {
		transAmountFirst = amount * pairRateValue[firstPairName][0]
		rateFirst = pairRateValue[firstPairName][0]
	} else {
		transAmountFirst = amount / pairRateValue[firstPairName][2]
		rateFirst = pairRateValue[firstPairName][2]
	}
	tmpAssetSecond := strings.TrimSuffix(firstPairNameLower, assetLower)

	//second steps
	if strings.HasPrefix(pairRate[0], tmpAssetSecond) {
		transAmountSecond = transAmountFirst * pairRateValue[pairRate[0]][0]
		rateSecond = pairRateValue[pairRate[0]][0]
	} else {
		transAmountSecond = transAmountFirst / pairRateValue[pairRate[0]][2]
		rateSecond = pairRateValue[pairRate[0]][2]
	}

	if strings.HasPrefix(pairRate[1], assetLower) {
		result = transAmountSecond / pairRateValue[pairRate[1]][2]
		rateThird = pairRateValue[pairRate[1]][2]
	} else {
		result = transAmountSecond * pairRateValue[pairRate[1]][0]
		rateThird = pairRateValue[pairRate[1]][0]
	}

	timeSince := time.Since(t)
	if result > amount {
		log.Printf("%d Asset: %s First Pair %s %.2f %.2f Secong Pair %s %.2f %.2f Third Pair %s %.2f %.2f - %v",
			count, asset, firstPairName, rateFirst, transAmountFirst, pairRate[0], rateSecond,
			transAmountSecond, pairRate[1], rateThird, result, timeSince)
	}
}

func GetPairFromJSONHuobi(asset string) map[string][]string {
	pair := ""
	switch asset {
	case "USDT":
		pair = fmt.Sprintf("data/datahuobi/%s/%s_pairH_pairH_pairH.json", asset, asset)
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
