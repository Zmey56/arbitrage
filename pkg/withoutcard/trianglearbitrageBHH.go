package withoutcard

import (
	"github.com/Zmey56/arbitrage/pkg/getinfobinance"
	"github.com/Zmey56/arbitrage/pkg/getinfohuobi"
	"log"
	"strings"
	"time"
)

func TriangleArbitrageBHH(asset string, amount float64) {
	start := time.Now()
	pair := getinfobinance.GetPairFromJSONBHH(asset)

	for firstPairName, secondThirdPairs := range pair {
		GetTriangleArbitrageBHH(start, asset, firstPairName, secondThirdPairs, amount)
	}
}

func GetTriangleArbitrageBHH(t time.Time, asset, firstPairName string, pairsArray []string, amount float64) {
	for count, pp := range pairsArray {
		CalculateTriangleArbitrageBHH(t, count, asset, firstPairName, pp, amount)

	}
}

func CalculateTriangleArbitrageBHH(t time.Time, count int, asset, firstPairName, pp string, amount float64) {
	pairRate := strings.Split(pp, "|")
	pairRateValue := getinfohuobi.GetRatePairTriangleHuobi(pairRate)              //rate for pairs from Huobi
	pairRateValueB := getinfobinance.GetRatePairTriangle([]string{firstPairName}) //rate for pair from Binance

	if pairRateValueB[firstPairName][0] == 0.0 || pairRateValue[pairRate[0]][0] == 0 || pairRateValue[pairRate[1]][0] == 0 ||
		pairRateValueB[firstPairName][2] == 0.0 || pairRateValue[pairRate[0]][2] == 0 || pairRateValue[pairRate[1]][2] == 0 {
		return
	}

	firstPairNameLower := strings.ToLower(firstPairName)
	assetLower := strings.ToLower(asset)
	var transAmountFirst, transAmountSecond, result float64
	var rateFirst, rateSecond, rateThird float64

	// first step
	if strings.HasPrefix(firstPairNameLower, assetLower) {
		transAmountFirst = amount * pairRateValueB[firstPairName][0]
		rateFirst = pairRateValueB[firstPairName][0]
	} else {
		transAmountFirst = amount / pairRateValueB[firstPairName][2]
		rateFirst = pairRateValueB[firstPairName][2]
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
