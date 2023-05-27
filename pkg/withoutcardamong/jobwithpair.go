package withoutcardamong

import (
	"github.com/Zmey56/arbitrage/pkg/getdatahuobi"
	"github.com/Zmey56/arbitrage/pkg/getdataokx"
	"github.com/Zmey56/arbitrage/pkg/getinfobinance"
	"log"
)

func TriangleArbitrageBetweenExc(fiat string) {

	pairB := getinfobinance.GetPairBinanceWithout(fiat)

	pairH := getdatahuobi.GetPairHuobiWithout(fiat)

	pairO := getdataokx.GetPairOKXWithout(fiat)

	comPairBH := CommonPair(pairB, pairH)
	comPairBO := CommonPair(pairB, pairO)
	comPairOH := CommonPair(pairO, pairH)

	log.Println(len(comPairBH))
	log.Println(len(comPairBO))
	log.Println(len(comPairOH))

	ratePairBH := getinfobinance.GetRatePairTriangleB(comPairBH)

	ratePairHB := getdatahuobi.GetRatePairTriangleH(comPairBH)

	// {bid, bidVolume, ask, askVolume}
	for pair, value := range ratePairBH {
		tmpResultFirst := ((ratePairHB[pair][0] - value[2]) / value[2]) * 100
		if tmpResultFirst > 1 && tmpResultFirst < 100 {
			log.Printf("Pair %s, first price - %.2f, second price - %.2f, profit - %.2f", pair, value[2],
				ratePairHB[pair][0], tmpResultFirst)
		}
	}

	log.Println("_____________________________________________________________________________")

	for pair, value := range ratePairBH {
		tmpResultSecond := ((value[0] - ratePairHB[pair][2]) / ratePairHB[pair][2]) * 100
		if tmpResultSecond > 1 && tmpResultSecond < 100 {
			log.Printf("Pair %s, first price - %.2f, second price - %.2f, profit - %.2f", pair, ratePairHB[pair][2],
				value[0], tmpResultSecond)
		}
	}

	log.Println("_____________________________________________________________________________")
	log.Println("_____________________________________________________________________________")

	ratePairBO := getinfobinance.GetRatePairTriangleB(comPairBO)

	ratePairOB := getdataokx.GetRatePairTriangleO(comPairBO, fiat)

	for pair, value := range ratePairBO {
		tmpResultFirst := ((ratePairOB[pair][0] - value[2]) / value[2]) * 100
		if tmpResultFirst > 1 && tmpResultFirst < 100 {
			log.Printf("Pair %s, first price - %.2f, second price - %.2f, profit - %.2f", pair, value[2],
				ratePairOB[pair][0], tmpResultFirst)
		}
	}

	log.Println("_____________________________________________________________________________")

	for pair, value := range ratePairBO {
		tmpResultSecond := ((value[0] - ratePairOB[pair][2]) / ratePairOB[pair][2]) * 100
		if tmpResultSecond > 1 && tmpResultSecond < 100 {
			log.Printf("Pair %s, first price - %.2f, second price - %.2f, profit - %.2f", pair, ratePairOB[pair][2],
				value[0], tmpResultSecond)
		}
	}

	log.Println("_____________________________________________________________________________")
	log.Println("_____________________________________________________________________________")

	ratePairHO := getdatahuobi.GetRatePairTriangleH(comPairOH)
	ratePairOH := getdataokx.GetRatePairTriangleO(comPairOH, fiat)

	for pair, value := range ratePairHO {
		tmpResultFirst := ((ratePairOH[pair][0] - value[2]) / value[2]) * 100
		if tmpResultFirst > 1 && tmpResultFirst < 100 {
			log.Printf("Pair %s, first price - %.2f, second price - %.2f, profit - %.2f", pair, value[2],
				ratePairOH[pair][0], tmpResultFirst)
		}
	}

	log.Println("_____________________________________________________________________________")

	for pair, value := range ratePairHO {
		tmpResultSecond := ((value[0] - ratePairOH[pair][2]) / ratePairOH[pair][2]) * 100
		if tmpResultSecond > 1 && tmpResultSecond < 100 {
			log.Printf("Pair %s, first price - %.2f, second price - %.2f, profit - %.2f", pair, ratePairOH[pair][2],
				value[0], tmpResultSecond)
		}
	}

}
