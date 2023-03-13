package p2pokx

import (
	"github.com/Zmey56/arbitrage/pkg/getdataokx"
	"github.com/Zmey56/arbitrage/pkg/getinfookx"
	"github.com/Zmey56/arbitrage/pkg/result"
	"log"
	"math"
	"strconv"
	"sync"
	"time"
)

func P2P2stepsOKX(fiat string, paramUser getinfookx.ParametersOKX) {

	//get all assets from OKX for this fiat
	currencyarr := getdataokx.GetCurrencyOKX(fiat)

	var wg sync.WaitGroup
	for _, a := range currencyarr {
		wg.Add(1)
		go func(a string) {
			defer wg.Done()
			getResultP2P2(a, fiat, paramUser)

		}(a)
	}
	wg.Wait()

}

func getResultP2P2(a, fiat string, paramUser getinfookx.ParametersOKX) {

	order_buy := getdataokx.GetDataP2POKXBuy(fiat, a, paramUser)

	if len(order_buy.Data.Sell) > 1 {
		var transAmountFloat float64
		if paramUser.Amount != "" {
			tmpTransAmountFloat, err := strconv.ParseFloat(paramUser.Amount, 64)
			if err != nil {
				log.Println("Can't convert transAmount", err)
			}
			transAmountFloat = tmpTransAmountFloat
		} else {
			tmpTransAmountFloat, err := strconv.ParseFloat(order_buy.Data.Sell[0].AvailableAmount, 64)
			if err != nil {
				log.Println("Can't convert MaxTradeLimit", err)
			}
			transAmountFloat = tmpTransAmountFloat
			paramUser.Amount = strconv.Itoa(int(transAmountFloat))
		}

		price_b, _ := strconv.ParseFloat(order_buy.Data.Sell[0].Price, 64)

		transAmountFirst := transAmountFloat / price_b
		//second step

		printResultP2P2(a, fiat, transAmountFirst, price_b, order_buy, paramUser)

	} else {
		log.Printf("Order buy is empty, fiat - %s, assets - %s, param %+v\n", fiat, a, paramUser)
	}

}

func printResultP2P2(a, fiat string, transAmountFirst, price_b float64,
	order_buy getdataokx.OKXBuy, paramUser getinfookx.ParametersOKX) {

	var assetSell = a

	//third steps

	order_sell := getdataokx.GetDataP2POKXSell(assetSell, fiat, paramUser)

	if len(order_sell.Data.Buy) < 2 {
		log.Printf("Order sell is empty, fiat - %s, assets - %s, param %+v\n", fiat, a, paramUser)
	} else {
		profitResult := deltaBuySell(order_buy, order_sell, a, fiat, paramUser)

		result.CheckResultSaveSend2Steps(profitResult, paramUser.Border)
	}
}

func deltaBuySell(ob getdataokx.OKXBuy, os getdataokx.OKXSell, asset, fiat string, pu getinfookx.ParametersOKX) result.ResultP2P2steps {
	res := result.ResultP2P2steps{}

	firstB, _ := strconv.ParseFloat(ob.Data.Sell[0].Price, 64)
	res.PriceB = firstB
	secondB, _ := strconv.ParseFloat(ob.Data.Sell[1].Price, 64)
	res.PriceBSecond = secondB

	firstS, _ := strconv.ParseFloat(os.Data.Buy[0].Price, 64)
	res.PriceS = firstS
	secondS, _ := strconv.ParseFloat(os.Data.Buy[1].Price, 64)
	res.PriceSSecond = secondS

	res.DeltaBuySell = (firstB - firstS) / firstB
	res.DeltaFirstSecondB = ((firstB - secondB) / firstB) * 100
	res.DeltaFirstSecondS = ((secondS - firstS) / firstS) * 100

	sumB := 0.0
	sumDeltaB := 0.0
	tmpB := 0.0

	for _, j := range ob.Data.Sell {
		tmpPP, _ := strconv.ParseFloat(j.Price, 64)
		sumDeltaB = sumDeltaB + (tmpPP - tmpB)
		tmpB = tmpPP
		sumB = sumB + tmpPP
		tmpVB, _ := strconv.ParseFloat(j.AvailableAmount, 64)
		if tmpVB > res.GiantVolB {
			res.GiantVolB = tmpVB
			res.GiantPriceB = tmpB
		}
	}

	meanB := sumB / float64(len(ob.Data))
	res.MeanPriceB = meanB

	res.DeltaGiantPriceB = res.GiantPriceB - res.PriceB

	varianceB := 0.0
	for _, valueB := range ob.Data.Sell {
		tmpPP2, _ := strconv.ParseFloat(valueB.Price, 64)
		diff := tmpPP2 - meanB
		varianceB += diff * diff
	}
	varianceB /= float64(len(ob.Data.Sell))
	res.SDPriceB = math.Sqrt(varianceB)

	sumS := 0.0
	sumDeltaS := 0.0
	tmpS := 0.0

	for _, i := range os.Data.Buy {
		tmpPS, _ := strconv.ParseFloat(i.Price, 64)
		sumDeltaS = sumDeltaS + (tmpPS - tmpS)
		tmpS = tmpPS
		sumS = sumS + tmpPS
		tmpVS, _ := strconv.ParseFloat(i.AvailableAmount, 64)
		if tmpVS > res.GiantVolS {
			res.GiantVolS = tmpVS
			res.GiantPriceS = tmpS
		}
	}

	meanS := sumS / float64(len(os.Data))
	res.MeanPriceS = meanS

	res.DeltaGiantPriceS = res.GiantPriceS - res.PriceS

	varianceS := 0.0
	for _, valueS := range os.Data.Buy {
		tmpPS2, _ := strconv.ParseFloat(valueS.Price, 64)
		diff := tmpPS2 - meanS
		varianceS += diff * diff
	}
	varianceS /= float64(len(os.Data.Buy))
	res.SDPriceS = math.Sqrt(varianceS)

	res.DeltaMean = meanB - meanS

	weightedSumB := 0.0
	for i := 0; i < len(ob.Data.Sell); i++ {
		tmp_w, _ := strconv.ParseFloat(ob.Data.Sell[i].AvailableAmount, 64)
		tmpPP3, _ := strconv.ParseFloat(ob.Data.Sell[i].Price, 64)
		weightedSumB += tmpPP3 * tmp_w
	}

	sumOfWeightsB := 0.0
	for i := 0; i < len(ob.Data.Sell); i++ {
		tmp_w_2, _ := strconv.ParseFloat(ob.Data.Sell[i].AvailableAmount, 64)
		sumOfWeightsB += tmp_w_2
	}

	res.MeanWeighB = weightedSumB / sumOfWeightsB

	weightedSumS := 0.0
	for j := 0; j < len(os.Data); j++ {
		tmp_ws, _ := strconv.ParseFloat(os.Data[j].TradeCount, 64)
		tmpPS3, _ := strconv.ParseFloat(os.Data[j].Price, 64)
		weightedSumS += tmpPS3 * tmp_ws
	}

	sumOfWeightsS := 0.0
	for i := 0; i < len(os.Data); i++ {
		tmp_ws_2, _ := strconv.ParseFloat(os.Data[i].TradeCount, 64)
		sumOfWeightsS += tmp_ws_2
	}

	res.MeanWeighS = weightedSumS / sumOfWeightsS

	res.DeltaMeanWeight = res.MeanWeighB - res.MeanWeighS

	res.AdvToalBuy = ob.TotalCount
	res.AdvToalSell = os.TotalCount
	res.DeltaADV = 100 * ((float64(res.AdvToalSell) - float64(res.AdvToalBuy)) / float64(res.AdvToalSell))

	res.FiatUnit = fiat
	res.Asset = asset
	res.Merchant = pu.IsMerchant == "true"
	res.DataTime = time.Now()
	res.MarketOne = "Huobi"
	res.MarketTwo = "Huobi"
	res.User.FirstUser = "Taker"
	res.User.SecondUser = "Taker"
	res.PaymentBuy = result.PaymentMetodsHuobi(ob)
	res.PaymentSell = result.PaymentMetodsHuobi(os)

	res.Amount, _ = strconv.ParseFloat(pu.Amount, 64)

	return res
}
