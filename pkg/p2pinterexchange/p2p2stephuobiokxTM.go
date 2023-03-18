package p2pinterexchange

import (
	"fmt"
	"github.com/Zmey56/arbitrage/pkg/commonfunction"
	"github.com/Zmey56/arbitrage/pkg/getdatahuobi"
	"github.com/Zmey56/arbitrage/pkg/getdataokx"
	"github.com/Zmey56/arbitrage/pkg/getinfohuobi"
	"github.com/Zmey56/arbitrage/pkg/result"
	"github.com/Zmey56/arbitrage/pkg/workinghuobi"
	"github.com/Zmey56/arbitrage/pkg/workingokx"
	"log"
	"math"
	"strconv"
	"sync"
	"time"
)

func P2P2stepsHuobiOKXTM(fiat string, paramUser getinfohuobi.ParametersHuobi) {
	//get all assets from binance for this fiat
	currencyarrH := getdatahuobi.GetCurrencyHuobi(fiat)
	log.Println("Huobi", fiat, " - ", currencyarrH)
	currencyarrO := getdataokx.GetCurrencyOKX(fiat)

	currencyarr := commonfunction.CommonElement(currencyarrO, currencyarrH)
	log.Println("OKX", fiat, " - ", currencyarrO)

	log.Println("Result", currencyarr)
	var wg sync.WaitGroup
	for _, a := range currencyarr {
		wg.Add(1)
		go func(a string) {
			defer wg.Done()
			getResultP2P2HOTM(a, fiat, paramUser)

		}(a)
	}
	wg.Wait()

}

func getResultP2P2HOTM(a, fiat string, paramUser getinfohuobi.ParametersHuobi) {
	coinidmap := workinghuobi.GetCoinIDHuobo(fiat)

	order_buy := getdatahuobi.GetDataP2PHuobi(coinidmap[a], coinidmap[fiat], "sell", paramUser)
	//log.Printf("%+v\n", order_buy)

	if len(order_buy.Data) > 1 {
		var transAmountFloat float64
		if paramUser.Amount != "" {
			tmpTransAmountFloat, err := strconv.ParseFloat(paramUser.Amount, 64)
			if err != nil {
				log.Println("Can't convert transAmount", err)
			}
			transAmountFloat = tmpTransAmountFloat
		} else {
			tmpTransAmountFloat, err := strconv.ParseFloat(order_buy.Data[0].MaxTradeLimit, 64)
			if err != nil {
				log.Println("Can't convert MaxTradeLimit", err)
			}
			transAmountFloat = tmpTransAmountFloat
			paramUser.Amount = strconv.Itoa(int(transAmountFloat))
		}

		price_b, _ := strconv.ParseFloat(order_buy.Data[0].Price, 64)

		transAmountFirst := transAmountFloat / price_b
		//second step

		printResultP2P2stepsHuobiOKXTM(a, fiat, transAmountFirst, price_b, order_buy, paramUser)

	} else {
		log.Printf("Order buy is empty, fiat - %s, assets - %s, param %+v\n", fiat, a, paramUser)
	}

}

func printResultP2P2stepsHuobiOKXTM(fiat, a string, transAmountFirst, price_b float64,
	order_buy getdatahuobi.Huobi, huobi getinfohuobi.ParametersHuobi) {
	var assetSell = a
	paramUserO := workingokx.GetParamOKX(fiat)
	paramUserO.Amount = fmt.Sprintf("%v", transAmountFirst)
	paramUserO.IsMerchant = strconv.FormatBool(huobi.IsMerchant == "true")

	//third steps
	log.Println(assetSell)

	order_sell := getdataokx.GetDataP2POKXBuy(assetSell, fiat, paramUserO)

	if len(order_sell.Data.Sell) < 2 {
		log.Printf("Order sell is empty, fiat - %s, assets - %s, param %+v\n", fiat, a, paramUserO)
	} else {
		profitResult := deltaBuySellHOTM(order_buy, order_sell, fiat, a, huobi)

		result.CheckResultSaveSend2Steps(profitResult, huobi.Border)
	}
}

func deltaBuySellHOTM(ob getdatahuobi.Huobi, os getdataokx.OKXBuy, asset, fiat string,
	pu getinfohuobi.ParametersHuobi) result.ResultP2P2steps {
	res := result.ResultP2P2steps{}

	firstB, _ := strconv.ParseFloat(ob.Data[0].Price, 64)
	res.PriceB = firstB
	secondB, _ := strconv.ParseFloat(ob.Data[1].Price, 64)
	res.PriceBSecond = secondB

	firstS, _ := strconv.ParseFloat(os.Data.Sell[0].Price, 64)
	res.PriceS = firstS
	secondS, _ := strconv.ParseFloat(os.Data.Sell[0].Price, 64)
	res.PriceSSecond = secondS

	res.DeltaBuySell = ((firstS - firstB) / firstB) * 100
	res.DeltaFirstSecondB = ((firstB - secondB) / firstB) * 100
	res.DeltaFirstSecondS = ((secondS - firstS) / firstS) * 100

	sumB := 0.0
	sumDeltaB := 0.0
	tmpB := 0.0

	for _, j := range ob.Data {
		tmpPP, _ := strconv.ParseFloat(j.Price, 64)
		sumDeltaB = sumDeltaB + (tmpPP - tmpB)
		tmpB = tmpPP
		sumB = sumB + tmpPP
		tmpVB, _ := strconv.ParseFloat(j.TradeCount, 64)
		if tmpVB > res.GiantVolB {
			res.GiantVolB = tmpVB
			res.GiantPriceB = tmpB
		}
	}

	meanB := sumB / float64(len(ob.Data))
	res.MeanPriceB = meanB

	res.DeltaGiantPriceB = ((res.PriceB - res.GiantPriceB) / res.PriceB) * 100

	varianceB := 0.0
	for _, valueB := range ob.Data {
		tmpPP2, _ := strconv.ParseFloat(valueB.Price, 64)
		diff := tmpPP2 - meanB
		varianceB += diff * diff
	}
	varianceB /= float64(len(ob.Data))
	res.SDPriceB = math.Sqrt(varianceB)

	sumS := 0.0
	sumDeltaS := 0.0
	tmpS := 0.0

	// Mean of sell adv
	for _, i := range os.Data.Sell {
		tmpSMean, _ := strconv.ParseFloat(i.Price, 64)
		sumDeltaS = sumDeltaS + (tmpSMean - tmpS)
		tmpS = tmpSMean
		sumS = sumS + tmpSMean
		tmpVS, _ := strconv.ParseFloat(i.AvailableAmount, 64)
		if tmpVS > res.GiantVolS {
			res.GiantVolS = tmpVS
			res.GiantPriceS = tmpS
		}
	}

	meanS := sumS / float64(len(os.Data.Sell))
	res.MeanPriceS = meanS

	res.DeltaGiantPriceS = ((res.PriceS - res.GiantPriceS) / res.PriceS) * 100

	varianceS := 0.0
	for _, valueS := range os.Data.Sell {
		valueStmp, _ := strconv.ParseFloat(valueS.Price, 64)
		diff := valueStmp - meanS
		//log.Println("TEST", diff)
		varianceS += diff * diff
	}
	varianceS /= float64(len(os.Data.Sell))
	res.SDPriceS = math.Sqrt(varianceS)

	res.DeltaMean = ((meanS - meanB) / meanB) * 100

	weightedSumB := 0.0
	for i := 0; i < len(ob.Data); i++ {
		tmp_w, _ := strconv.ParseFloat(ob.Data[i].TradeCount, 64)
		tmpPP3, _ := strconv.ParseFloat(ob.Data[i].Price, 64)
		weightedSumB += tmpPP3 * tmp_w
	}

	sumOfWeightsB := 0.0
	for i := 0; i < len(ob.Data); i++ {
		tmp_w_2, _ := strconv.ParseFloat(ob.Data[i].TradeCount, 64)
		sumOfWeightsB += tmp_w_2
	}

	res.MeanWeighB = weightedSumB / sumOfWeightsB

	weightedSumS := 0.0
	for j := 0; j < len(os.Data.Sell); j++ {
		tmp_ws, _ := strconv.ParseFloat(os.Data.Sell[j].AvailableAmount, 64)
		tmpPrice, _ := strconv.ParseFloat(os.Data.Sell[j].Price, 64)
		weightedSumS += tmpPrice * tmp_ws
	}

	sumOfWeightsS := 0.0
	for i := 0; i < len(os.Data.Sell); i++ {
		tmp_ws_2, _ := strconv.ParseFloat(os.Data.Sell[i].AvailableAmount, 64)
		sumOfWeightsS += tmp_ws_2
	}

	res.MeanWeighS = weightedSumS / sumOfWeightsS

	res.DeltaMeanWeight = ((res.MeanWeighS - res.MeanWeighB) / res.MeanWeighB) * 100

	res.AdvToalBuy = ob.TotalCount
	res.AdvToalSell = os.Data.Total
	res.DeltaADV = 100 * ((float64(res.AdvToalSell) - float64(res.AdvToalBuy)) / float64(res.AdvToalSell))
	//log.Println("TEST", float64(res.AdvToalBuy), "-", float64(res.AdvToalSell), (float64(res.AdvToalBuy)-float64(res.AdvToalSell))/float64(res.AdvToalSell))
	res.DeltaGiant = ((res.GiantPriceS - res.GiantPriceB) / res.GiantPriceB) * 100

	res.FiatUnit = fiat
	res.Asset = asset
	res.Merchant = pu.IsMerchant == "true"
	res.DataTime = time.Now()
	res.MarketOne = "Huobi"
	res.MarketTwo = "OKX"
	res.User.FirstUser = "Taker"
	res.User.SecondUser = "Maker"

	res.PaymentBuy = result.PaymentMetodsHuobi(ob)
	res.PaymentSell = os.Data.Sell[0].PaymentMethods

	res.Amount, _ = strconv.ParseFloat(pu.Amount, 64)

	return res
}
