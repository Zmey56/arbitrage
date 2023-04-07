package p2pbinance

import (
	"github.com/Zmey56/arbitrage/pkg/getdata"
	"github.com/Zmey56/arbitrage/pkg/getinfobinance"
	"github.com/Zmey56/arbitrage/pkg/workingbinance"
	"log"
	"strconv"
	"strings"
	"sync"
)

func P2S3stepsTSS(fiat string, paramUser workingbinance.ParametersBinance) {

	//get all assets and pair from binance for this fiat

	pair := getinfobinance.GetPairFromJSONPairPair(fiat)

	//get information about orders with binance
	var wg sync.WaitGroup
	for a, _ := range pair {
		wg.Add(1)
		go func(a string) {
			defer wg.Done()
			GetResultP2S3TSS(a, fiat, pair, paramUser)
		}(a)
	}
	wg.Wait()
}

func GetResultP2S3TSS(a, fiat string, pair map[string][]string, paramUser workingbinance.ParametersBinance) {

	//first step
	order_buy := getdata.GetDataP2PBinance(a, fiat, "Buy", paramUser)
	var transAmountFloat float64
	if paramUser.TransAmount != "" {
		tmpTransAmountFloat, err := strconv.ParseFloat(paramUser.TransAmount, 64)
		if err != nil {
			log.Println("Can't convert transAmount", err)
		}
		transAmountFloat = tmpTransAmountFloat
	} else {
		tmpTransAmountFloat, err := strconv.ParseFloat(order_buy.Data[0].Adv.DynamicMaxSingleTransAmount, 64)
		if err != nil {
			log.Println("Can't convert dynamicMaxSingleTransAmount", err)
		}
		transAmountFloat = tmpTransAmountFloat
		paramUser.TransAmount = strconv.Itoa(int(transAmountFloat))
	}

	if len(order_buy.Data) > 0 {

		price_b := order_buy.Data[0].Adv.Price

		transAmountFirst := transAmountFloat / price_b

		//second step
		//for each element in arrange get gorutine

		var wg sync.WaitGroup
		for _, p := range pair[a] {
			wg.Add(1)

			go func(p string) {
				defer wg.Done()
				PrintResultP2S3TSS(a, p, fiat, transAmountFirst, price_b,
					order_buy, paramUser)
			}(string(p))

		}
		wg.Wait()
	} else {
		log.Printf("Order buy is empty, fiat - %s, assets - %s, param %+v\n", fiat, a, paramUser)
	}

}

func PrintResultP2S3TSS(a, p, fiat string, transAmountFirst, price_b float64,
	order_buy getinfobinance.AdvertiserAdv, paramUser workingbinance.ParametersBinance) {

	pairArray := strings.Split(p, "|")

	if len(pairArray) < 2 {
		log.Println("Not correct value in array of currency")
		return
	}

	pairRate := getinfobinance.GetRatePair(pairArray)
	secondPair := strings.TrimSpace(pairArray[0])
	thirdPair := strings.TrimSpace(pairArray[1])

	if pairRate[secondPair] == 0.0 || pairRate[thirdPair] == 0.0 {
		return
	}

	var transAmountSecond float64

	log.Println("Asset", a, "Second Pair", secondPair, "Result", strings.HasPrefix(strings.TrimSpace(a), secondPair))

	if strings.HasPrefix(secondPair, a) {
		transAmountSecond = (transAmountFirst * pairRate[secondPair])
	} else {
		transAmountSecond = (transAmountFirst / pairRate[secondPair])
	}
	//third steps
	var transAmountThird float64

	if !strings.HasPrefix(thirdPair, fiat) {
		transAmountThird = (transAmountSecond * pairRate[thirdPair])
	} else {
		transAmountThird = (transAmountSecond / pairRate[thirdPair])
	}

	if transAmountThird > 10000 {
		log.Printf("Fiat %s - %s, Asset, %s - %.8f(%.2f), Second pair %s - %.8f(%.8f), Third pair %s - %.8f, RESULT: %.8f",
			fiat, paramUser.TransAmount, a, price_b, transAmountFirst, secondPair, pairRate[secondPair], transAmountSecond, thirdPair,
			pairRate[thirdPair], transAmountThird)
	}
}
