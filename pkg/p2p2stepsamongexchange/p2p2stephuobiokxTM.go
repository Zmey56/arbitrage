package p2p2stepsamongexchange

import (
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
	assetsO := getdataokx.GetCurrencyOKX(fiat)
	assetsH := getdatahuobi.GetCurrencyHuobi(fiat)
	//assetsB := make([]string, 0, len(assets))
	//for k, _ := range assets {
	//	assetsB = append(assetsB, k)
	//}
	assetsSymbol := commonfunction.CommonElement(assetsO, assetsH)
	//log.Println("assetsB", assetsB)
	//log.Println("assetsH", assetsH)
	//log.Println("assetsSymbol", assetsSymbol)

	var wg sync.WaitGroup
	for _, a := range assetsSymbol {
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

		printResultP2P2HOTM(a, fiat, transAmountFirst, price_b, order_buy, paramUser)

	} else {
		log.Printf("Order buy is empty, fiat - %s, assets - %s, param %+v\n", fiat, a, paramUser)
	}

}

func printResultP2P2HOTM(a, fiat string, transAmountFirst, price_b float64,
	order_buy getdatahuobi.Huobi, paramUser getinfohuobi.ParametersHuobi) {

	paramUserO := workingokx.GetParamOKX(fiat)
	paramUserO.Amount = paramUser.Amount
	if paramUser.IsMerchant == "true" {
		paramUserO.IsMerchant = "true"
	} else {
		paramUserO.IsMerchant = "false"
	}

	order_sell := getdataokx.GetDataP2POKXBuy(fiat, a, paramUserO)

	if len(order_sell.Data.Sell) < 2 {
		log.Printf("Order sell is empty, fiat - %s, assets - %s, param %+v\n", fiat, a, paramUserO)
	} else {
		//log.Println("fiat", fiat, "Asset", a, "LEN")

		profitResult := deltaBuySellHOTM(order_buy, order_sell, a, fiat, paramUser)
		result.CheckResultSaveSend2Steps(profitResult, paramUser.Border)
	}
}

func deltaBuySellHOTM(ob getdatahuobi.Huobi, os getdataokx.OKXBuy, asset, fiat string, pu getinfohuobi.ParametersHuobi) result.ResultP2P2steps {
	res := result.ResultP2P2steps{}

	firstB, _ := strconv.ParseFloat(ob.Data[0].Price, 64)
	res.PriceB = firstB
	secondB, _ := strconv.ParseFloat(ob.Data[1].Price, 64)
	res.PriceBSecond = secondB

	firstS, _ := strconv.ParseFloat(os.Data.Sell[0].Price, 64)
	res.PriceS = firstS
	secondS, _ := strconv.ParseFloat(os.Data.Sell[0].Price, 64)
	res.PriceSSecond = secondS

	res.DeltaBuySell = 100 * (firstS - firstB) / firstB
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

	res.DeltaGiantPriceB = 100 * (res.PriceB - res.GiantPriceB) / res.PriceB

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

	res.DeltaGiantPriceS = 100 * (res.PriceS - res.GiantPriceS) / res.PriceS

	varianceS := 0.0
	for _, valueS := range os.Data.Sell {
		valueStmp, _ := strconv.ParseFloat(valueS.Price, 64)
		diff := valueStmp - meanS
		//log.Println("TEST", diff)
		varianceS += diff * diff
	}
	varianceS /= float64(len(os.Data.Sell))
	res.SDPriceS = math.Sqrt(varianceS)

	res.DeltaMean = 100 * (meanS - meanB) / meanB

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

	res.MeanWeight = (weightedSumB + weightedSumS) / (sumOfWeightsB + sumOfWeightsS)

	res.DeltaMeanWeight = 100 * (res.MeanWeighS - res.MeanWeighB) / res.MeanWeighB

	res.AdvToalBuy = ob.TotalCount
	res.AdvToalSell = os.Data.Total
	//res.DeltaADV = 100 * ((float64(res.AdvToalBuy) - float64(res.AdvToalSell)) / float64(res.AdvToalBuy))
	res.DeltaADV = 100 * ((float64(res.AdvToalSell) - float64(res.AdvToalBuy)) / float64(res.AdvToalSell))
	res.DeltaGiant = 100.0 * (res.GiantPriceS - res.GiantPriceB) / res.GiantPriceB
	//log.Println("TEST", float64(res.AdvToalBuy), "-", float64(res.AdvToalSell), (float64(res.AdvToalBuy)-float64(res.AdvToalSell))/float64(res.AdvToalSell))

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

	res.DeltaSDb = (res.SDPriceB / res.PriceB) * 100
	res.DeltaSDs = (res.SDPriceS / res.PriceS) * 100
	res.DeltaSD = ((res.SDPriceS - res.SDPriceB) / res.SDPriceB) * 100

	res.Amount, _ = strconv.ParseFloat(pu.Amount, 64)

	return res
}
