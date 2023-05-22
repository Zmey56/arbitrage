package p2p2stepsamongexchange

import (
	"fmt"
	"github.com/Zmey56/arbitrage/pkg/commonfunction"
	"github.com/Zmey56/arbitrage/pkg/getdata"
	"github.com/Zmey56/arbitrage/pkg/getdatahuobi"
	"github.com/Zmey56/arbitrage/pkg/getinfobinance"
	"github.com/Zmey56/arbitrage/pkg/result"
	"github.com/Zmey56/arbitrage/pkg/workingbinance"
	"github.com/Zmey56/arbitrage/pkg/workinghuobi"
	"log"
	"math"
	"strconv"
	"strings"
	"sync"
	"time"
)

func P2P2stepsBinanceHuobiTT(fiat string, paramUser workingbinance.ParametersBinance) {
	//get all assets from binance and huobi for this fiat
	assetsH := getdatahuobi.GetCurrencyHuobi(fiat)
	assetsB := getdata.GetAssetsLocalBinance(fiat)
	assetsSymbol := commonfunction.CommonElement(assetsB, assetsH)

	var wg sync.WaitGroup
	for _, a := range assetsSymbol {
		wg.Add(1)
		go func(asset string) {
			getResultP2P2stepsBinanceHuobiTT(fiat, asset, paramUser)
			wg.Done()
		}(a)
	}
	wg.Wait()

}

func getResultP2P2stepsBinanceHuobiTT(fiat, a string, paramUser workingbinance.ParametersBinance) {
	//log.Println("START SECOND FUNCTION", fiat, " - ", a)
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
		//log.Println("Second")
		price_b := order_buy.Data[0].Adv.Price

		transAmountFirst := transAmountFloat / price_b
		//second step

		printResultP2P2stepsBinanceHuobiTT(fiat, a, transAmountFirst, price_b, order_buy, paramUser)
	}
}

func printResultP2P2stepsBinanceHuobiTT(fiat, a string, transAmountFirst, price_b float64,
	order_buy getinfobinance.AdvertiserAdv, binance workingbinance.ParametersBinance) {
	coinidmap := workinghuobi.GetCoinIDHuobo(fiat)
	paramUserH := workinghuobi.GetParamHuobi(fiat)
	paramUserH.Amount = binance.TransAmount
	if binance.PublisherType != "merchant" {
		paramUserH.IsMerchant = "false"
	} else {
		paramUserH.IsMerchant = "true"
	}

	//third steps
	var assetSell = a

	//third steps

	if coinidmap[strings.ToUpper(assetSell)] != 0 {
		order_sell := getdatahuobi.GetDataP2PHuobi(coinidmap[strings.ToUpper(assetSell)], coinidmap[fiat],
			"buy", paramUserH)
		if len(order_sell.Data) < 2 {
			log.Printf("Order sell is empty, fiat - %s, assets - %s, param %+v\n", fiat, a, paramUserH)
		} else {
			profitResult := result.ResultP2P{}

			price_s, _ := strconv.ParseFloat(order_sell.Data[0].Price, 64)

			transAmountThird := price_s * transAmountFirst

			transAmountFloat, err := strconv.ParseFloat(binance.TransAmount, 64)
			if err != nil {
				log.Printf("Problem with convert transAmount to float, err - %v", err)
			}
			profitResult.Amount = binance.TransAmount
			profitResult.Market.First = "Binance"
			profitResult.Merchant.FirstMerch = (binance.PublisherType == "merchant")
			profitResult.User.FirstUser = "Taker"
			profitResult.Market.Second = ""
			profitResult.Market.Third = "Huobi"
			profitResult.Merchant.ThirdMerch = (paramUserH.IsMerchant == "true")
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
			profitResult.PaymentSell = result.PaymentMetodsHuobi(order_sell)
			profitResult.LinkAssetsSell = fmt.Sprintf("https://www.huobi.com/en-us/fiat-crypto/trade/sell-%s-%s/", assetSell, strings.ToLower(fiat))
			profitResult.ProfitValue = transAmountThird - transAmountFloat
			profitResult.ProfitPercet = (((transAmountThird - transAmountFloat) / transAmountFloat) * 100)
			profitResult.TotalAdvBuy = order_buy.Total
			profitResult.TotalAdvSell = order_sell.TotalCount
			profitResult.AdvNoBuy = order_buy.Data[0].Adv.AdvNo
			profitResult.AdvNoSell = strconv.Itoa(order_sell.Data[0].UID)

			result.CheckResultSaveSend2Steps(profitResult, binance.Border)

		}
	}
}

func deltaBuySellBHTT(ob getinfobinance.AdvertiserAdv, os getdatahuobi.Huobi, asset, fiat string,
	pu workingbinance.ParametersBinance) result.ResultP2P2steps {
	res := result.ResultP2P2steps{}
	tmpData := []float64{}
	tmpDataW := []float64{}

	firstB := ob.Data[0].Adv.Price
	res.PriceB = firstB
	secondB := ob.Data[1].Adv.Price
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

	res.DeltaGiantPriceS = ((res.PriceS - res.GiantPriceS) / res.PriceS) * 100

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

	res.AdvToalBuy = ob.Total
	res.AdvToalSell = os.TotalCount
	res.DeltaADV = 100 * ((float64(res.AdvToalSell) - float64(res.AdvToalBuy)) / float64(res.AdvToalSell))
	res.DeltaGiant = 100.0 * (res.GiantPriceS - res.GiantPriceB) / res.GiantPriceB

	res.FiatUnit = fiat
	res.Asset = asset
	res.Merchant = pu.PublisherType == "merchant"
	res.DataTime = time.Now()
	res.MarketOne = "Binance"
	res.MarketTwo = "Huobi"
	res.User.FirstUser = "Taker"
	res.User.SecondUser = "Taker"
	res.PaymentBuy = result.PaymentMetods(ob)
	res.PaymentSell = result.PaymentMetodsHuobi(os)
	res.DeltaSDb = (res.SDPriceB / res.PriceB) * 100
	res.DeltaSDs = (res.SDPriceS / res.PriceS) * 100
	res.DeltaSD = ((res.SDPriceS - res.SDPriceB) / res.SDPriceB) * 100

	res.Amount, _ = strconv.ParseFloat(pu.TransAmount, 64)

	log.Println("tmpData", tmpData)
	log.Println("tmpDataW", tmpDataW)
	res.MeanWeightSD = commonfunction.WeightedStandardDeviation(tmpData, tmpDataW)
	res.DeltaWSD = (res.MeanWeightSD / res.PriceB) * 100

	return res
}
