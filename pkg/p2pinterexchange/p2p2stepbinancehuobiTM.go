package p2pinterexchange

import (
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

func P2P2stepsBinanceHuobiTM(fiat string, paramUser workingbinance.ParametersBinance) {
	//get all assets from binance for this fiat

	assets := getdata.GetAssets(fiat)
	assets_symbol := make([]string, 0, len(assets))
	//assets_name := make([]string, 0, len(assets))

	for k, _ := range assets {
		assets_symbol = append(assets_symbol, k)
		//assets_name = append(assets_name, v)
	}

	var wg sync.WaitGroup
	for _, a := range assets_symbol {
		wg.Add(1)
		go func(asset string) {
			getResultP2P2stepsBinanceHuobiTM(fiat, asset, paramUser)
			wg.Done()
		}(a)
	}
	wg.Wait()

}

func getResultP2P2stepsBinanceHuobiTM(fiat, a string, paramUser workingbinance.ParametersBinance) {
	//log.Println("START SECOND FUNCTION", fiat, " - ", a)
	order_buy := getdata.GetDataP2PBinance(a, fiat, "Buy", paramUser)
	//log.Println("ORDER_BUY", order_buy)
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

		printResultP2P2stepsBinanceHuobiTM(fiat, a, transAmountFirst, price_b, order_buy, paramUser)

	} else {
		log.Printf("Order buy is empty, fiat - %s, assets - %s, param %+v\n", fiat, a, paramUser)
	}
}

func printResultP2P2stepsBinanceHuobiTM(fiat, a string, transAmountFirst, price_b float64,
	order_buy getinfobinance.AdvertiserAdv, binance workingbinance.ParametersBinance) {
	coinidmap := workinghuobi.GetCoinIDHuobo(fiat)
	paramUserH := workinghuobi.GetParamHuobi(fiat)

	//third steps
	var assetSell = a

	//third steps

	if coinidmap[strings.ToUpper(assetSell)] != 0 {
		order_sell := getdatahuobi.GetDataP2PHuobi(coinidmap[strings.ToUpper(assetSell)], coinidmap[fiat],
			"sell", paramUserH)
		//log.Printf("%+v\n\n", order_sell)
		if len(order_sell.Data) < 2 {
			log.Printf("Order sell is empty, fiat - %s, assets - %s, param %+v\n", fiat, a, paramUserH)
		} else {
			//log.Println("fiat", fiat, "Asset", a, "LEN")

			profitResult := deltaBuySellBHTM(order_buy, order_sell, a, fiat, binance)
			//log.Println("PROFIT", profitResult)
			//log.Println("binance.Border", binance.Border, len(order_sell.Data), "-", len(order_buy.Data))
			result.CheckResultSaveSend2Steps(profitResult, binance.Border)
		}
	}
}

func deltaBuySellBHTM(ob getinfobinance.AdvertiserAdv, os getdatahuobi.Huobi, asset, fiat string,
	pu workingbinance.ParametersBinance) result.ResultP2P2steps {
	res := result.ResultP2P2steps{}

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
		tmpSMean, _ := strconv.ParseFloat(i.Price, 64)
		sumDeltaS = sumDeltaS + (tmpSMean - tmpS)
		tmpS = tmpSMean
		sumS = sumS + tmpSMean
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
		tmp_w, _ := strconv.ParseFloat(ob.Data[i].Adv.SurplusAmount, 64)
		weightedSumB += ob.Data[i].Adv.Price * tmp_w
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
		tmpPrice, _ := strconv.ParseFloat(os.Data[j].Price, 64)
		weightedSumS += tmpPrice * tmp_ws
	}

	sumOfWeightsS := 0.0
	for i := 0; i < len(os.Data); i++ {
		tmp_ws_2, _ := strconv.ParseFloat(os.Data[i].TradeCount, 64)
		sumOfWeightsS += tmp_ws_2
	}

	res.MeanWeighS = weightedSumS / sumOfWeightsS

	res.DeltaMeanWeight = ((res.MeanWeighS - res.MeanWeighB) / res.MeanWeighB) * 100

	res.AdvToalBuy = ob.Total
	res.AdvToalSell = os.TotalCount
	res.DeltaADV = 100 * ((float64(res.AdvToalSell) - float64(res.AdvToalBuy)) / float64(res.AdvToalSell))
	//log.Println("TEST", float64(res.AdvToalBuy), "-", float64(res.AdvToalSell), (float64(res.AdvToalBuy)-float64(res.AdvToalSell))/float64(res.AdvToalSell))
	res.DeltaGiant = ((res.GiantPriceS - res.GiantPriceB) / res.GiantPriceB) * 100

	res.FiatUnit = fiat
	res.Asset = asset
	res.Merchant = pu.PublisherType == "merchant"
	res.DataTime = time.Now()
	res.MarketOne = "Binance"
	res.MarketTwo = "Huobi"
	res.User.FirstUser = "Taker"
	res.User.SecondUser = "Maker"

	res.PaymentBuy = result.PaymentMetods(ob)
	res.PaymentSell = result.PaymentMetodsHuobi(os)

	res.Amount, _ = strconv.ParseFloat(pu.TransAmount, 64)

	return res
}
