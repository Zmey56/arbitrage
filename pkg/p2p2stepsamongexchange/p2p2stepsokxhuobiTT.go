package p2p2stepsamongexchange

import (
	"github.com/Zmey56/arbitrage/pkg/commonfunction"
	"github.com/Zmey56/arbitrage/pkg/getdatahuobi"
	"github.com/Zmey56/arbitrage/pkg/getdataokx"
	"github.com/Zmey56/arbitrage/pkg/getinfookx"
	"github.com/Zmey56/arbitrage/pkg/result"
	"github.com/Zmey56/arbitrage/pkg/workinghuobi"
	"log"
	"math"
	"strconv"
	"strings"
	"sync"
	"time"
)

func P2P2stepsOKXHuobiTT(fiat string, paramUser getinfookx.ParametersOKX) {

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
			getResultP2P2stepsOKXHuobiTT(a, fiat, paramUser)

		}(a)
	}
	wg.Wait()

}

func getResultP2P2stepsOKXHuobiTT(a, fiat string, paramUser getinfookx.ParametersOKX) {

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

		printResultP2P2stepsOKXHuobiTT(a, fiat, transAmountFirst, price_b, order_buy, paramUser)

	} else {
		log.Printf("Order buy is empty, fiat - %s, assets - %s, param %+v\n", fiat, a, paramUser)
	}

}

func printResultP2P2stepsOKXHuobiTT(a, fiat string, transAmountFirst, price_b float64,
	order_buy getdataokx.OKXBuy, paramUser getinfookx.ParametersOKX) {

	coinidmap := workinghuobi.GetCoinIDHuobo(fiat)
	paramUserH := workinghuobi.GetParamHuobi(fiat)
	paramUserH.Amount = paramUser.Amount
	if paramUser.IsMerchant != "true" {
		paramUserH.IsMerchant = "false"
	} else {
		paramUserH.IsMerchant = "true"
	}

	//third steps
	var assetSell = a

	//third steps
	//log.Println(coinidmap, "assetSell", assetSell)

	if coinidmap[strings.ToUpper(assetSell)] != 0 {
		order_sell := getdatahuobi.GetDataP2PHuobi(coinidmap[strings.ToUpper(assetSell)], coinidmap[fiat],
			"buy", paramUserH)
		//log.Printf("len %v %+v\n\n", len(order_sell.Data), order_sell)
		if len(order_sell.Data) < 2 {
			log.Printf("Order sell is empty, fiat - %s, assets - %s, param %+v\n", fiat, a, paramUserH)
		} else {
			log.Println("paramUserH", paramUserH)

			profitResult := deltaBuySellOHTT(order_buy, order_sell, a, fiat, paramUser)
			result.CheckResultSaveSend2Steps(profitResult, paramUser.Border)
		}
	}
}

func deltaBuySellOHTT(ob getdataokx.OKXBuy, os getdatahuobi.Huobi, asset, fiat string, pu getinfookx.ParametersOKX) result.ResultP2P2steps {
	res := result.ResultP2P2steps{}
	tmpData := []float64{}
	tmpDataW := []float64{}

	firstB, _ := strconv.ParseFloat(ob.Data.Sell[0].Price, 64)
	res.PriceB = firstB
	secondB, _ := strconv.ParseFloat(ob.Data.Sell[1].Price, 64)
	res.PriceBSecond = secondB

	firstS, _ := strconv.ParseFloat(os.Data[0].Price, 64)
	res.PriceS = firstS
	secondS, _ := strconv.ParseFloat(os.Data[1].Price, 64)
	res.PriceSSecond = secondS

	res.DeltaBuySell = ((firstS - firstB) / firstB) * 100
	res.DeltaFirstSecondB = ((firstB - secondB) / firstB) * 100
	res.DeltaFirstSecondS = ((secondS - firstS) / firstS) * 100

	sumB := 0.0
	sumDeltaB := 0.0
	tmpB := 0.0

	for _, j := range ob.Data.Sell {
		tmpPP, _ := strconv.ParseFloat(j.Price, 64)
		tmpData = append(tmpData, tmpPP) //for weight SD
		sumDeltaB = sumDeltaB + (tmpPP - tmpB)
		tmpB = tmpPP
		sumB = sumB + tmpPP
		tmpVB, _ := strconv.ParseFloat(j.AvailableAmount, 64)
		if tmpVB > res.GiantVolB {
			res.GiantVolB = tmpVB
			res.GiantPriceB = tmpB
		}
	}

	meanB := sumB / float64(len(ob.Data.Sell))
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

	for _, i := range os.Data {
		tmpPS, _ := strconv.ParseFloat(i.Price, 64)
		tmpData = append(tmpData, tmpPS) //for weight SD
		sumDeltaS = sumDeltaS + (tmpPS - tmpS)
		tmpS = tmpPS
		sumS = sumS + tmpPS
		tmpVS, _ := strconv.ParseFloat(i.TradeCount, 64)
		if tmpVS > res.GiantVolS {
			res.GiantVolS = tmpVS
			res.GiantPriceS = tmpS
		}
	}

	meanS := sumS / float64(len(os.Data))
	res.MeanPriceS = meanS

	res.DeltaGiantPriceS = res.GiantPriceS - res.PriceS

	varianceS := 0.0
	for _, valueS := range os.Data {
		valueStmp, _ := strconv.ParseFloat(valueS.Price, 64)
		diff := valueStmp - meanS
		//log.Println("TEST", diff)
		varianceS += diff * diff
	}
	varianceS /= float64(len(os.Data))
	res.SDPriceS = math.Sqrt(varianceS)

	res.DeltaMean = ((meanS - meanB) / meanB) * 100

	weightedSumB := 0.0
	for i := 0; i < len(ob.Data.Sell); i++ {
		tmp_wb, _ := strconv.ParseFloat(ob.Data.Sell[i].AvailableAmount, 64)
		tmpDataW = append(tmpDataW, tmp_wb) //for weight SD
		tmpPP3, _ := strconv.ParseFloat(ob.Data.Sell[i].Price, 64)
		weightedSumB += tmpPP3 * tmp_wb
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
		tmpDataW = append(tmpDataW, tmp_ws) //for weight SD
		tmpPrice, _ := strconv.ParseFloat(os.Data[j].Price, 64)
		weightedSumS += tmpPrice * tmp_ws
	}

	sumOfWeightsS := 0.0
	for i := 0; i < len(os.Data); i++ {
		tmp_ws_2, _ := strconv.ParseFloat(os.Data[i].TradeCount, 64)
		sumOfWeightsS += tmp_ws_2
	}

	res.MeanWeighS = weightedSumS / sumOfWeightsS

	res.MeanWeight = (weightedSumB + weightedSumS) / (sumOfWeightsB + sumOfWeightsS)

	res.DeltaMeanWeight = ((res.MeanWeighS - res.MeanWeighB) / res.MeanWeighB) * 100

	res.AdvToalBuy = ob.Data.Total
	res.AdvToalSell = os.TotalCount
	res.DeltaADV = 100 * ((float64(res.AdvToalSell) - float64(res.AdvToalBuy)) / float64(res.AdvToalSell))

	res.FiatUnit = fiat
	res.Asset = asset
	res.Merchant = pu.IsMerchant == "true"
	res.DataTime = time.Now()
	res.MarketOne = "OKX"
	res.MarketTwo = "Huobi"
	res.User.FirstUser = "Taker"
	res.User.SecondUser = "Taker"
	res.PaymentBuy = ob.Data.Sell[0].PaymentMethods
	res.PaymentSell = result.PaymentMetodsHuobi(os)
	res.DeltaGiant = ((res.GiantPriceS - res.GiantPriceB) / res.GiantPriceB) * 100

	res.DeltaSDb = (res.SDPriceB / res.PriceB) * 100
	res.DeltaSDs = (res.SDPriceS / res.PriceS) * 100
	res.DeltaSD = ((res.SDPriceS - res.SDPriceB) / res.SDPriceB) * 100

	res.Amount, _ = strconv.ParseFloat(pu.Amount, 64)

	log.Println("tmpData", tmpData)
	log.Println("tmpDataW", tmpDataW)
	res.MeanWeightSD = commonfunction.WeightedStandardDeviation(tmpData, tmpDataW)
	res.DeltaWSD = (res.MeanWeightSD / res.PriceB) * 100

	return res
}
