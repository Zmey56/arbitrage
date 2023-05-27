package p2pbinance

import (
	"fmt"
	"github.com/Zmey56/arbitrage/pkg/getdata"
	"github.com/Zmey56/arbitrage/pkg/getinfobinance"
	"github.com/Zmey56/arbitrage/pkg/result"
	"github.com/Zmey56/arbitrage/pkg/workingbinance"
	"log"
	"strconv"
	"strings"
	"sync"
	"time"
)

func P2P3stepsTakerMaker(fiat string, paramUser workingbinance.ParametersBinance) {

	//get all assets and pair from binance for this fiat

	pair := getinfobinance.GetPairFromJSON(fiat)

	//get information about orders with binance
	var wg sync.WaitGroup
	for a, _ := range pair {
		wg.Add(1)
		go func(a string) {
			defer wg.Done()
			GetResultP2P3TM(a, fiat, pair, paramUser)
		}(a)
	}
	wg.Wait()
}

func GetResultP2P3TM(a, fiat string, pair map[string][]string, paramUser workingbinance.ParametersBinance) {

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

		pair_rate := getinfobinance.GetRatePair(pair[a])

		var wg sync.WaitGroup
		for p := range pair_rate {
			wg.Add(1)

			go func(p string) {
				defer wg.Done()
				PrintResultP2P3TM(p, a, fiat, transAmountFirst, price_b,
					pair_rate, order_buy, paramUser)
			}(p)

		}
		wg.Wait()
	}

}

func PrintResultP2P3TM(p, a, fiat string, transAmountFirst, price_b float64, pair_rate map[string]float64,
	order_buy getinfobinance.AdvertiserAdv, paramUser workingbinance.ParametersBinance) {

	profitResult := result.ResultP2P{}

	var transAmountSecond float64
	var assetSell string
	if strings.HasPrefix(p, a) {
		transAmountSecond = (transAmountFirst * pair_rate[p])
		assetSell = p[len(a):]
	} else {
		transAmountSecond = (transAmountFirst / pair_rate[p])
		assetSell = p[:(len(p) - len(a))]
	}
	//third steps
	order_sell := getdata.GetDataP2PBinance(assetSell, fiat,
		"Buy", paramUser)
	if len(order_sell.Data) > 0 {
		price_s := order_sell.Data[0].Adv.Price

		transAmountThird := price_s * transAmountSecond

		transAmountFloat, err := strconv.ParseFloat(paramUser.TransAmount, 64)
		if err != nil {
			log.Printf("Problem with convert transAmount to float, err - %v", err)
		}
		profitResult.Amount = paramUser.TransAmount
		profitResult.Market.First = "Binance"
		profitResult.Merchant.FirstMerch = (paramUser.PublisherType == "merchant")
		profitResult.User.FirstUser = "Taker"
		profitResult.Market.Second = "Binance"
		profitResult.Market.Third = "Binance"
		profitResult.User.ThirdUser = "Maker"
		profitResult.Merchant.ThirdMerch = (paramUser.PublisherType == "merchant")
		profitResult.Profit = transAmountThird > transAmountFloat
		profitResult.DataTime = time.Now()
		profitResult.Fiat = fiat
		profitResult.AssetsBuy = a
		profitResult.PriceAssetsBuy = price_b
		profitResult.PaymentBuy = result.PaymentMetods(order_buy)
		profitResult.LinkAssetsBuy = fmt.Sprintf("https://p2p.binance.com/en/trade/all-payments/%v?fiat=%v", a, fiat)
		profitResult.Pair = p
		profitResult.PricePair = pair_rate[p]
		profitResult.LinkMarket = result.ReturnLinkMarket(a, p)
		profitResult.AssetsSell = assetSell
		profitResult.PriceAssetsSell = price_s
		profitResult.PaymentSell = result.PaymentMetods(order_sell)
		profitResult.LinkAssetsSell = fmt.Sprintf("https://p2p.binance.com/en/trade/all-payments/%v?fiat=%v",
			assetSell, fiat)
		profitResult.ProfitValue = transAmountThird - transAmountFloat
		profitResult.ProfitPercet = (((transAmountThird - transAmountFloat) / transAmountFloat) * 100)
		profitResult.TotalAdvBuy = order_buy.Total
		profitResult.TotalAdvSell = order_sell.Total
		profitResult.AdvNoBuy = order_buy.Data[0].Adv.AdvNo
		profitResult.AdvNoSell = order_sell.Data[0].Adv.AdvNo

		//result2steps := getResultP2P2stepsBinance(fiat, a, assetSell, paramUser)
		result.CheckResultSaveSend(profitResult.User.FirstUser, profitResult.User.ThirdUser, paramUser.Border, paramUser.PercentUser, profitResult)

	}
}

//func getResultP2P2stepsBinance(fiat, assOne, assTwo string, paramUser workingbinance.ParametersBinance) result.ResultP2P2steps {
//	order_buy := getdata.GetDataP2PBinance(assOne, fiat, "Buy", paramUser)
//	transAmountFloat, err := strconv.ParseFloat(paramUser.TransAmount, 64)
//	if err != nil {
//		log.Println("Can't convert transAmount", err)
//	}
//
//	if len(order_buy.Data) > 1 {
//		price_b := order_buy.Data[0].Adv.Price
//
//		transAmountFirst := transAmountFloat / price_b
//		//second step
//
//		return printResultP2P2stepsBinance(fiat, assOne, assTwo, transAmountFirst, price_b, order_buy, paramUser)
//
//	} else {
//		log.Printf("Order buy is empty, fiat - %s, assets - %s, param %+v\n", fiat, assOne, paramUser)
//		return result.ResultP2P2steps{}
//	}
//}
//
//func printResultP2P2stepsBinance(fiat, a, b string, transAmountFirst, price_b float64,
//	order_buy getinfobinance.AdvertiserAdv, paramUser workingbinance.ParametersBinance) result.ResultP2P2steps {
//
//	//third steps
//	order_sell := getdata.GetDataP2PBinance(b, fiat, "Buy", paramUser)
//	if len(order_sell.Data) < 2 {
//		log.Printf("Order sell is empty, fiat - %s, assets - %s, param %+v\n", fiat, a, paramUser)
//		return result.ResultP2P2steps{}
//	} else {
//		profitResult := deltaBuySellBinance(order_buy, order_sell, a, fiat, paramUser)
//		return profitResult
//	}
//}
//
//func deltaBuySellBinance(ob, os getinfobinance.AdvertiserAdv, asset, fiat string, pu workingbinance.ParametersBinance) result.ResultP2P2steps {
//	res := result.ResultP2P2steps{}
//	//variable for find weight SD
//	tmpData := []float64{}
//	tmpDataW := []float64{}
//
//	firstB := ob.Data[0].Adv.Price
//	res.PriceB = firstB
//	secondB := ob.Data[1].Adv.Price
//	res.PriceBSecond = secondB
//
//	firstS := os.Data[0].Adv.Price
//	res.PriceS = firstS
//	secondS := os.Data[1].Adv.Price
//	res.PriceSSecond = secondS
//
//	res.DeltaBuySell = ((firstS - firstB) / firstB) * 100
//	res.DeltaFirstSecondB = ((firstB - secondB) / firstB) * 100
//	res.DeltaFirstSecondS = ((secondS - firstS) / firstS) * 100
//
//	sumB := 0.0
//	sumDeltaB := 0.0
//	tmpB := 0.0
//
//	for _, j := range ob.Data {
//		sumDeltaB = sumDeltaB + (j.Adv.Price - tmpB)
//		tmpB = j.Adv.Price
//		tmpData = append(tmpData, tmpB) //for weight SD
//		sumB = sumB + j.Adv.Price
//		tmpVB, _ := strconv.ParseFloat(j.Adv.SurplusAmount, 64)
//		if tmpVB > res.GiantVolB {
//			res.GiantVolB = tmpVB
//			res.GiantPriceB = tmpB
//		}
//	}
//
//	meanB := sumB / float64(len(ob.Data))
//	res.MeanPriceB = meanB
//
//	res.DeltaGiantPriceB = ((res.PriceB - res.GiantPriceB) / res.PriceB) * 100
//
//	varianceB := 0.0
//	for _, valueB := range ob.Data {
//		diff := valueB.Adv.Price - meanB
//		varianceB += diff * diff
//	}
//	varianceB /= float64(len(ob.Data))
//	res.SDPriceB = math.Sqrt(varianceB)
//
//	sumS := 0.0
//	sumDeltaS := 0.0
//	tmpS := 0.0
//
//	for _, i := range os.Data {
//		sumDeltaS = sumDeltaS + (i.Adv.Price - tmpS)
//		tmpS = i.Adv.Price
//		tmpData = append(tmpData, tmpS) //for weight SD
//		sumS = sumS + i.Adv.Price
//		tmpVS, _ := strconv.ParseFloat(i.Adv.SurplusAmount, 64)
//		if tmpVS > res.GiantVolS {
//			res.GiantVolS = tmpVS
//			res.GiantPriceS = tmpS
//		}
//	}
//
//	meanS := sumS / float64(len(os.Data))
//	res.MeanPriceS = meanS
//
//	res.DeltaGiantPriceS = ((res.GiantPriceS - res.PriceS) / res.PriceS) * 100
//
//	varianceS := 0.0
//	for _, valueS := range os.Data {
//		diff := valueS.Adv.Price - meanS
//		varianceS += diff * diff
//	}
//	varianceS /= float64(len(os.Data))
//	res.SDPriceS = math.Sqrt(varianceS)
//
//	res.DeltaMean = ((meanS - meanB) / meanB) * 100
//
//	weightedSumB := 0.0
//	for i := 0; i < len(ob.Data); i++ {
//		tmp_w, _ := strconv.ParseFloat(ob.Data[i].Adv.SurplusAmount, 64)
//		tmpDataW = append(tmpDataW, tmp_w) //for weight SD
//		weightedSumB += ob.Data[i].Adv.Price * tmp_w
//	}
//
//	sumOfWeightsB := 0.0
//	for i := 0; i < len(ob.Data); i++ {
//		tmp_w_2, _ := strconv.ParseFloat(ob.Data[i].Adv.SurplusAmount, 64)
//		sumOfWeightsB += tmp_w_2
//	}
//
//	res.MeanWeighB = weightedSumB / sumOfWeightsB
//
//	weightedSumS := 0.0
//	for j := 0; j < len(os.Data); j++ {
//		tmp_ws, _ := strconv.ParseFloat(os.Data[j].Adv.SurplusAmount, 64)
//		tmpDataW = append(tmpDataW, tmp_ws) //for weight SD
//		weightedSumS += os.Data[j].Adv.Price * tmp_ws
//	}
//
//	sumOfWeightsS := 0.0
//	for i := 0; i < len(os.Data); i++ {
//		tmp_ws_2, _ := strconv.ParseFloat(os.Data[i].Adv.SurplusAmount, 64)
//		sumOfWeightsS += tmp_ws_2
//	}
//
//	res.MeanWeighS = weightedSumS / sumOfWeightsS
//	res.MeanWeight = (weightedSumB + weightedSumS) / (sumOfWeightsB + sumOfWeightsS)
//
//	res.DeltaMeanWeight = ((res.MeanWeighS - res.MeanWeighB) / res.MeanWeighB) * 100
//	res.MeanWeightSD = commonfunction.WeightedStandardDeviation(tmpData, tmpDataW)
//	res.DeltaWSD = (res.MeanWeightSD / res.PriceB) * 100
//
//	res.AdvToalBuy = ob.Total
//	res.AdvToalSell = os.Total
//	res.DeltaADV = 100 * ((float64(res.AdvToalSell) - float64(res.AdvToalBuy)) / float64(res.AdvToalSell))
//	res.DeltaGiant = ((res.GiantPriceS - res.GiantPriceB) / res.GiantPriceB) * 100
//
//	res.FiatUnit = fiat
//	res.Asset = asset
//	res.Merchant = (pu.PublisherType == "merchant")
//	res.DataTime = time.Now()
//	res.MarketOne = "Binance"
//	res.MarketTwo = "Binance"
//	res.User.FirstUser = "Taker"
//	res.User.SecondUser = "Taker"
//	res.PaymentBuy = result.PaymentMetods(ob)
//	res.PaymentSell = result.PaymentMetods(os)
//	res.DeltaSDb = (res.SDPriceB / res.PriceB) * 100
//	res.DeltaSDs = (res.SDPriceS / res.PriceS) * 100
//	res.DeltaSD = ((res.SDPriceS - res.SDPriceB) / res.SDPriceB) * 100
//
//	res.Amount, _ = strconv.ParseFloat(pu.TransAmount, 64)
//
//	return res
//}
