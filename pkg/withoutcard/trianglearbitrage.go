package withoutcard

import (
	"github.com/Zmey56/arbitrage/pkg/getinfobinance"
	"log"
	"strings"
)

func TriangleArbitrage(asset string, amount float64) {
	pair := getinfobinance.GetPairFromJSONPairPairPair(asset)

	for firstPairName, secondThirdPairs := range pair {
		GetTriangleArbitrage(asset, firstPairName, secondThirdPairs, amount)
	}
}

func GetTriangleArbitrage(asset, firstPairName string, pairsArray []string, amount float64) {
	for _, pp := range pairsArray {
		CalculateTriangleArbitrage(asset, firstPairName, pp, amount)
	}
}

func CalculateTriangleArbitrage(asset, firstPairName, pp string, amount float64) {
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

	if result > amount {
		log.Printf("Asset: %s First Pair %s %.2f %.2f Secong Pair %s %.2f %.2f Third Pair %s %.2f %.2f",
			asset, firstPairName, rateFirst, transAmountFirst, pairRate[0], rateSecond,
			transAmountSecond, pairRate[1], rateThird, result)
	}
}
