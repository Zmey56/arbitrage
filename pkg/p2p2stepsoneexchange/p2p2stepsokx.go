package p2p2stepsoneexchange

import (
	"fmt"
	"github.com/Zmey56/arbitrage/pkg/commonfunction"
	"github.com/Zmey56/arbitrage/pkg/getdataokx"
	"github.com/Zmey56/arbitrage/pkg/getinfookx"
	"github.com/Zmey56/arbitrage/pkg/result"
	"log"
	"math"
	"strconv"
	"strings"
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
			getResultP2P2OKX(a, fiat, paramUser)

		}(a)
	}
	wg.Wait()

}

func getResultP2P2OKX(a, fiat string, paramUser getinfookx.ParametersOKX) {

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

		printResultP2P2OKX(a, fiat, transAmountFirst, price_b, order_buy, paramUser)

	} else {
		log.Printf("Order buy is empty, fiat - %s, assets - %s, param %+v\n", fiat, a, paramUser)
	}

}

func printResultP2P2OKX(a, fiat string, transAmountFirst, price_b float64,
	order_buy getdataokx.OKXBuy, paramUser getinfookx.ParametersOKX) {

	var assetSell = a

	//third steps
	log.Println(assetSell)

	order_sell := getdataokx.GetDataP2POKXSell(fiat, assetSell, paramUser)

	if len(order_sell.Data.Buy) < 2 {
		log.Printf("Order sell is empty, fiat - %s, assets - %s, param %+v\n", fiat, a, paramUser)
	} else {
		price_s, _ := strconv.ParseFloat(order_sell.Data.Buy[0].Price, 64)
		transAmountFloat, err := strconv.ParseFloat(paramUser.Amount, 64)
		if err != nil {
			log.Printf("Problem with convert transAmount to float, err - %v", err)
		}

		transAmountThird := price_s * transAmountFirst

		profitResult := result.ResultP2P{}
		profitResult.Amount = paramUser.Amount
		profitResult.Market.First = "OKX"
		profitResult.Merchant.FirstMerch = (paramUser.IsMerchant == "true")
		profitResult.User.FirstUser = "Taker"
		profitResult.Market.Second = ""
		profitResult.Market.Third = "OKX"
		profitResult.Merchant.ThirdMerch = (paramUser.IsMerchant == "true")
		profitResult.User.ThirdUser = "Taker"
		profitResult.Profit = transAmountThird > transAmountFloat
		profitResult.DataTime = time.Now()
		profitResult.Fiat = fiat
		profitResult.AssetsBuy = a
		profitResult.PriceAssetsBuy = price_b
		profitResult.PaymentBuy = order_buy.Data.Sell[0].PaymentMethods
		profitResult.LinkAssetsBuy = fmt.Sprintf("https://www.okx.com/p2p-markets/%s/buy-%s/", strings.ToLower(fiat), strings.ToLower(a))
		profitResult.AssetsSell = assetSell
		profitResult.PriceAssetsSell = price_s
		profitResult.PaymentSell = order_sell.Data.Buy[0].PaymentMethods
		profitResult.LinkAssetsSell = fmt.Sprintf("https://www.okx.com/p2p-markets/%s/sell-%s", fiat, assetSell)
		profitResult.ProfitValue = transAmountThird - transAmountFloat
		profitResult.ProfitPercet = (((transAmountThird - transAmountFloat) / transAmountFloat) * 100)
		profitResult.TotalAdvBuy = order_buy.Data.Total
		profitResult.TotalAdvSell = order_sell.Data.Total
		profitResult.AdvNoBuy = order_buy.Data.Sell[0].ID
		profitResult.AdvNoSell = order_sell.Data.Buy[0].ID

		result.CheckResultSaveSend2Steps(profitResult, paramUser.Border)
	}
}

// may be sometime
func deltaBuySellOKX(ob getdataokx.OKXBuy, os getdataokx.OKXSell, asset, fiat string, pu getinfookx.ParametersOKX) result.ResultP2P2steps {
	res := result.ResultP2P2steps{}
	tmpData := []float64{}  // for Weighted mean Buy and Sell
	tmpDataW := []float64{} // for Weighted mean Buy and Sell

	firstB, _ := strconv.ParseFloat(ob.Data.Sell[0].Price, 64)
	res.PriceB = firstB
	secondB, _ := strconv.ParseFloat(ob.Data.Sell[1].Price, 64)
	res.PriceBSecond = secondB

	firstS, _ := strconv.ParseFloat(os.Data.Buy[0].Price, 64)
	res.PriceS = firstS
	secondS, _ := strconv.ParseFloat(os.Data.Buy[1].Price, 64)
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

	res.DeltaGiantPriceB = ((res.PriceB - res.GiantPriceB) / res.PriceB) * 100

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

	meanS := sumS / float64(len(os.Data.Buy))
	res.MeanPriceS = meanS

	res.DeltaGiantPriceS = ((res.GiantPriceS - res.PriceS) / res.PriceS) * 100

	varianceS := 0.0
	for _, valueS := range os.Data.Buy {
		tmpPS2, _ := strconv.ParseFloat(valueS.Price, 64)
		diff := tmpPS2 - meanS
		varianceS += diff * diff
	}
	varianceS /= float64(len(os.Data.Buy))
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
	for j := 0; j < len(os.Data.Buy); j++ {
		tmp_ws, _ := strconv.ParseFloat(os.Data.Buy[j].AvailableAmount, 64)
		tmpDataW = append(tmpDataW, tmp_ws) //for weight SD
		tmpPS3, _ := strconv.ParseFloat(os.Data.Buy[j].Price, 64)
		weightedSumS += tmpPS3 * tmp_ws
	}

	sumOfWeightsS := 0.0
	for i := 0; i < len(os.Data.Buy); i++ {
		tmp_ws_2, _ := strconv.ParseFloat(os.Data.Buy[i].AvailableAmount, 64)
		sumOfWeightsS += tmp_ws_2
	}

	res.MeanWeighS = weightedSumS / sumOfWeightsS

	res.MeanWeight = (weightedSumB + weightedSumS) / (sumOfWeightsB + sumOfWeightsS)

	res.DeltaMeanWeight = ((res.MeanWeighS - res.MeanWeighB) / res.MeanWeighB) * 100

	res.AdvToalBuy = ob.Data.Total
	res.AdvToalSell = os.Data.Total
	res.DeltaADV = 100 * ((float64(res.AdvToalSell) - float64(res.AdvToalBuy)) / float64(res.AdvToalSell))

	res.FiatUnit = fiat
	res.Asset = asset
	res.Merchant = pu.IsMerchant == "true"
	res.DataTime = time.Now()
	res.MarketOne = "OKX"
	res.MarketTwo = "OKX"
	res.User.FirstUser = "Taker"
	res.User.SecondUser = "Taker"
	res.PaymentBuy = ob.Data.Sell[0].PaymentMethods
	res.PaymentSell = os.Data.Buy[0].PaymentMethods
	res.DeltaGiant = ((res.GiantPriceS - res.GiantPriceB) / res.GiantPriceB) * 100
	res.DeltaSDb = (res.SDPriceB / res.PriceB) * 100
	res.DeltaSDs = (res.SDPriceS / res.PriceS) * 100
	res.DeltaSD = ((res.SDPriceS - res.SDPriceB) / res.SDPriceB) * 100
	//res.DeltaBuySell = ((firstS - firstB) / firstB) * 100

	res.Amount, _ = strconv.ParseFloat(pu.Amount, 64)

	res.MeanWeightSD = commonfunction.WeightedStandardDeviation(tmpData, tmpDataW)
	res.DeltaWSD = (res.MeanWeightSD / res.PriceB) * 100

	return res
}
