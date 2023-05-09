package withoutcard

import (
	"github.com/Zmey56/arbitrage/pkg/getinfobinance"
	"log"
	"strings"
	"time"
)

func TriangleArbitrageBinance(asset string, amount float64) {
	start := time.Now()
	pair := getinfobinance.GetPairFromJSONPairPairPair(asset)

	for firstPairName, secondThirdPairs := range pair {
		GetTriangleArbitrage(start, asset, firstPairName, secondThirdPairs, amount)
	}
}

func GetTriangleArbitrage(t time.Time, asset, firstPairName string, pairsArray []string, amount float64) {
	for count, pp := range pairsArray {
		CalculateTriangleArbitrage(t, count, asset, firstPairName, pp, amount)

	}
}

func CalculateTriangleArbitrage(t time.Time, count int, asset, firstPairName, pp string, amount float64) {
	pairRate := strings.Split(pp, "|")
	pairRate = append(pairRate, firstPairName)

	pairRateValue := getinfobinance.GetRatePairTriangle(pairRate)

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
		log.Printf("%d Asset: %s First Pair %s %.2f %.2f Secong Pair %s %.2f %.2f Third Pair %s %.2f %.2f - %v",
			count, asset, firstPairName, rateFirst, transAmountFirst, pairRate[0], rateSecond,
			transAmountSecond, pairRate[1], rateThird, result, timeSince)
	}
}
