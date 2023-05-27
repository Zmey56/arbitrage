package p2p2stepsamongexchange

import (
	"fmt"
	"github.com/Zmey56/arbitrage/pkg/commonfunction"
	"github.com/Zmey56/arbitrage/pkg/getdata"
	"github.com/Zmey56/arbitrage/pkg/getdataokx"
	"github.com/Zmey56/arbitrage/pkg/getinfobinance"
	"github.com/Zmey56/arbitrage/pkg/result"
	"github.com/Zmey56/arbitrage/pkg/workingbinance"
	"github.com/Zmey56/arbitrage/pkg/workingokx"
	"log"
	"math"
	"strconv"
	"sync"
	"time"
)

func P2P2stepsBinanceOKXTT(fiat string, paramUser workingbinance.ParametersBinance) {
	//get all assets from binance for this fiat

	assetsB := getdata.GetAssetsLocalBinance(fiat)
	assetsO := getdataokx.GetCurrencyOKX(fiat)

	assetsSymbol := commonfunction.CommonElement(assetsB, assetsO)

	var wg sync.WaitGroup
	for _, a := range assetsSymbol {
		wg.Add(1)
		go func(asset string) {
			getResultP2P2stepsBinanceOKXTT(fiat, asset, paramUser)
			wg.Done()
		}(a)
	}
	wg.Wait()

}

func getResultP2P2stepsBinanceOKXTT(fiat, a string, paramUser workingbinance.ParametersBinance) {

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

	if len(order_buy.Data) > 1 {
		price_b := order_buy.Data[0].Adv.Price

		transAmountFirst := transAmountFloat / price_b
		//second step

		printResultP2P2stepsBinanceOKXTT(fiat, a, transAmountFirst, price_b, order_buy, paramUser)

	}
}

func printResultP2P2stepsBinanceOKXTT(fiat, a string, transAmountFirst, price_b float64,
	order_buy getinfobinance.AdvertiserAdv, binance workingbinance.ParametersBinance) {
	var assetSell = a
	paramUserO := workingokx.GetParamOKX(fiat)
	paramUserO.Amount = binance.TransAmount
	paramUserO.IsMerchant = strconv.FormatBool(binance.PublisherType == "merchant")

	//third steps

	order_sell := getdataokx.GetDataP2POKXSell(fiat, assetSell, paramUserO)

	if len(order_sell.Data.Buy) > 1 {
		price_s, _ := strconv.ParseFloat(order_sell.Data.Buy[0].Price, 64)
		transAmountFloat, err := strconv.ParseFloat(binance.TransAmount, 64)
		if err != nil {
			log.Printf("Problem with convert transAmount to float, err - %v", err)
		}

		transAmountThird := price_s * transAmountFirst

		profitResult := result.ResultP2P{}
		profitResult.Amount = binance.TransAmount
		profitResult.Market.First = "Binance"
		profitResult.Merchant.FirstMerch = (binance.PublisherType == "merchant")
		profitResult.User.FirstUser = "Taker"
		profitResult.Market.Second = ""
		profitResult.Market.Third = "OKX"
		profitResult.Merchant.ThirdMerch = (binance.PublisherType == "merchant")
		profitResult.User.ThirdUser = "Taker"
		profitResult.Profit = transAmountThird > transAmountFloat
		profitResult.DataTime = time.Now()
		profitResult.Fiat = fiat
		profitResult.AssetsBuy = a
		profitResult.PriceAssetsBuy = price_b
		profitResult.PaymentBuy = result.PaymentMetods(order_buy)
		profitResult.LinkAssetsBuy = fmt.Sprintf("https://p2p.binance.com/en/trade/all-payments/%v?fiat=%v", a, fiat)
		profitResult.AssetsSell = assetSell
		profitResult.PriceAssetsSell = price_s
		profitResult.PaymentSell = order_sell.Data.Buy[0].PaymentMethods
		profitResult.LinkAssetsSell = fmt.Sprintf("https://www.okx.com/p2p-markets/%s/sell-%s", fiat, assetSell)
		profitResult.ProfitValue = transAmountThird - transAmountFloat
		profitResult.ProfitPercet = (((transAmountThird - transAmountFloat) / transAmountFloat) * 100)
		profitResult.TotalAdvBuy = order_buy.Total
		profitResult.TotalAdvSell = order_sell.Data.Total
		profitResult.AdvNoBuy = order_buy.Data[0].Adv.AdvNo
		profitResult.AdvNoSell = order_sell.Data.Buy[0].ID

		result.CheckResultSaveSend2Steps(binance.Border, binance.PercentUser, profitResult)
	}
}

func deltaBuySellBOTT(ob getinfobinance.AdvertiserAdv, os getdataokx.OKXSell, asset, fiat string,
	pu workingbinance.ParametersBinance) result.ResultP2P2steps {
	res := result.ResultP2P2steps{}
	tmpData := []float64{}
	tmpDataW := []float64{}

	firstB := ob.Data[0].Adv.Price
	res.PriceB = firstB
	secondB := ob.Data[1].Adv.Price
	res.PriceBSecond = secondB

	firstS, _ := strconv.ParseFloat(os.Data.Buy[0].Price, 64)
	res.PriceS = firstS
	secondS, _ := strconv.ParseFloat(os.Data.Buy[0].Price, 64)
	res.PriceSSecond = secondS

	res.DeltaBuySell = ((firstS - firstB) / firstB) * 100
	res.DeltaFirstSecondB = ((firstB - secondB) / firstB) * 100
	res.DeltaFirstSecondS = ((secondS - firstS) / firstS) * 100

	sumB := 0.0
	sumDeltaB := 0.0
	tmpB := 0.0

	for _, j := range ob.Data {
		sumDeltaB = sumDeltaB + (j.Adv.Price - tmpB)
		tmpB = j.Adv.Price
		tmpData = append(tmpData, tmpB) //for weight SD
		sumB = sumB + j.Adv.Price
		tmpVB, _ := strconv.ParseFloat(j.Adv.SurplusAmount, 64)
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
		diff := valueB.Adv.Price - meanB
		varianceB += diff * diff
	}
	varianceB /= float64(len(ob.Data))
	res.SDPriceB = math.Sqrt(varianceB)

	sumS := 0.0
	sumDeltaS := 0.0
	tmpS := 0.0

	// Mean of sell adv
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

	res.DeltaGiantPriceS = ((res.PriceS - res.GiantPriceS) / res.PriceS) * 100

	varianceS := 0.0
	for _, valueS := range os.Data.Buy {
		valueStmp, _ := strconv.ParseFloat(valueS.Price, 64)
		diff := valueStmp - meanS
		varianceS += diff * diff
	}
	varianceS /= float64(len(os.Data.Buy))
	res.SDPriceS = math.Sqrt(varianceS)

	res.DeltaMean = ((meanS - meanB) / meanB) * 100

	weightedSumB := 0.0
	for i := 0; i < len(ob.Data); i++ {
		tmp_wb, _ := strconv.ParseFloat(ob.Data[i].Adv.SurplusAmount, 64)
		tmpDataW = append(tmpDataW, tmp_wb) //for weight SD
		weightedSumB += ob.Data[i].Adv.Price * tmp_wb
	}

	sumOfWeightsB := 0.0
	for i := 0; i < len(ob.Data); i++ {
		tmp_w_2, _ := strconv.ParseFloat(ob.Data[i].Adv.SurplusAmount, 64)
		sumOfWeightsB += tmp_w_2
	}

	res.MeanWeighB = weightedSumB / sumOfWeightsB

	weightedSumS := 0.0
	for j := 0; j < len(os.Data.Buy); j++ {
		tmp_ws, _ := strconv.ParseFloat(os.Data.Buy[j].AvailableAmount, 64)
		tmpDataW = append(tmpDataW, tmp_ws) //for weight SD
		tmpPrice, _ := strconv.ParseFloat(os.Data.Buy[j].Price, 64)
		weightedSumS += tmpPrice * tmp_ws
	}

	sumOfWeightsS := 0.0
	for i := 0; i < len(os.Data.Buy); i++ {
		tmp_ws_2, _ := strconv.ParseFloat(os.Data.Buy[i].AvailableAmount, 64)
		sumOfWeightsS += tmp_ws_2
	}

	res.MeanWeighS = weightedSumS / sumOfWeightsS

	res.MeanWeight = (weightedSumB + weightedSumS) / (sumOfWeightsB + sumOfWeightsS)

	res.DeltaMeanWeight = ((res.MeanWeighS - res.MeanWeighB) / res.MeanWeighB) * 100

	res.AdvToalBuy = ob.Total
	res.AdvToalSell = os.Data.Total
	res.DeltaADV = 100 * ((float64(res.AdvToalSell) - float64(res.AdvToalBuy)) / float64(res.AdvToalSell))
	res.DeltaGiant = ((res.GiantPriceS - res.GiantPriceB) / res.GiantPriceB) * 100

	res.FiatUnit = fiat
	res.Asset = asset
	res.Merchant = pu.PublisherType == "merchant"
	res.DataTime = time.Now()
	res.MarketOne = "Binance"
	res.MarketTwo = "OKX"
	res.User.FirstUser = "Taker"
	res.User.SecondUser = "Taker"

	res.PaymentBuy = result.PaymentMetods(ob)
	res.PaymentSell = os.Data.Buy[0].PaymentMethods

	res.Amount, _ = strconv.ParseFloat(pu.TransAmount, 64)

	log.Println("tmpData", tmpData)
	log.Println("tmpDataW", tmpDataW)
	res.MeanWeightSD = commonfunction.WeightedStandardDeviation(tmpData, tmpDataW)
	res.DeltaWSD = (res.MeanWeightSD / res.PriceB) * 100

	return res
}
