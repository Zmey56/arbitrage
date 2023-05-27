package p2p2stepsamongexchange

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
	"strings"
	"sync"
	"time"
)

func P2P2stepsHuobiOKXTM(fiat string, paramUser getinfohuobi.ParametersHuobi) {
	assetsO := getdataokx.GetCurrencyOKX(fiat)
	assetsH := getdatahuobi.GetCurrencyHuobi(fiat)

	assetsSymbol := commonfunction.CommonElement(assetsO, assetsH)

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

	if len(order_sell.Data.Sell) > 1 {

		price_s, _ := strconv.ParseFloat(order_sell.Data.Sell[0].Price, 64)
		transAmountFloat, err := strconv.ParseFloat(paramUser.Amount, 64)
		if err != nil {
			log.Printf("Problem with convert transAmount to float, err - %v", err)
		}

		transAmountThird := price_s * transAmountFirst

		profitResult := result.ResultP2P{}
		profitResult.Amount = paramUser.Amount
		profitResult.Market.First = "Huobi"
		profitResult.Merchant.FirstMerch = (paramUser.IsMerchant == "true")
		profitResult.User.FirstUser = "Taker"
		profitResult.Market.Second = ""
		profitResult.Market.Third = "OKX"
		profitResult.Merchant.ThirdMerch = (paramUserO.IsMerchant == "true")
		profitResult.User.ThirdUser = "Maker"
		profitResult.Profit = transAmountThird > transAmountFloat
		profitResult.DataTime = time.Now()
		profitResult.Fiat = fiat
		profitResult.AssetsBuy = a
		profitResult.PriceAssetsBuy = price_b
		profitResult.PaymentBuy = result.PaymentMetodsHuobi(order_buy)
		profitResult.LinkAssetsBuy = fmt.Sprintf("https://www.huobi.com/en-us/fiat-crypto/trade/buy-%s-%s/", strings.ToLower(a), strings.ToLower(fiat))
		profitResult.AssetsSell = a
		profitResult.PriceAssetsSell = price_s
		profitResult.PaymentSell = order_sell.Data.Sell[0].PaymentMethods
		profitResult.LinkAssetsSell = fmt.Sprintf("https://www.okx.com/p2p-markets/%s/buy-%s/", strings.ToLower(fiat), strings.ToLower(a))
		profitResult.ProfitValue = transAmountThird - transAmountFloat
		profitResult.ProfitPercet = (((transAmountThird - transAmountFloat) / transAmountFloat) * 100)
		profitResult.TotalAdvBuy = order_buy.TotalCount
		profitResult.TotalAdvSell = order_sell.Data.Total
		profitResult.AdvNoBuy = strconv.Itoa(order_buy.Data[0].UID)
		profitResult.AdvNoSell = order_sell.Data.Sell[0].ID

		result.CheckResultSaveSend2Steps(paramUser.Border, paramUser.PercentUser, profitResult)
	}
}

func deltaBuySellHOTM(ob getdatahuobi.Huobi, os getdataokx.OKXBuy, asset, fiat string, pu getinfohuobi.ParametersHuobi) result.ResultP2P2steps {
	res := result.ResultP2P2steps{}
	tmpData := []float64{}
	tmpDataW := []float64{}

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
		tmpData = append(tmpData, tmpPP) //for weight SD
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
		tmpPS, _ := strconv.ParseFloat(i.Price, 64)
		tmpData = append(tmpData, tmpPS) //for weight SD
		sumDeltaS = sumDeltaS + (tmpPS - tmpS)
		tmpS = tmpPS
		sumS = sumS + tmpPS
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
		tmp_wb, _ := strconv.ParseFloat(ob.Data[i].TradeCount, 64)
		tmpDataW = append(tmpDataW, tmp_wb) //for weight SD
		tmpPP3, _ := strconv.ParseFloat(ob.Data[i].Price, 64)
		weightedSumB += tmpPP3 * tmp_wb
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
		tmpDataW = append(tmpDataW, tmp_ws) //for weight SD
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

	log.Println("tmpData", tmpData)
	log.Println("tmpDataW", tmpDataW)
	res.MeanWeightSD = commonfunction.WeightedStandardDeviation(tmpData, tmpDataW)
	res.DeltaWSD = (res.MeanWeightSD / res.PriceB) * 100

	return res
}
