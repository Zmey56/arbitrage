package Interexchange

import (
	"fmt"
	"github.com/Zmey56/arbitrage/pkg/getdata"
	"github.com/Zmey56/arbitrage/pkg/getdatahuobi"
	"github.com/Zmey56/arbitrage/pkg/getinfobinance"
	"github.com/Zmey56/arbitrage/pkg/getinfohuobi"
	"github.com/Zmey56/arbitrage/pkg/result"
	"github.com/Zmey56/arbitrage/pkg/workingbinance"
	"github.com/Zmey56/arbitrage/pkg/workinghuobi"
	"log"
	"strconv"
	"strings"
	"sync"
	"time"
)

func P2P3stepsTMBHH(fiat string, huobi getinfohuobi.ParametersHuobi,
	binance workingbinance.ParametersBinance) {
	//log.Println(huobi)
	allOrders := [][]result.ResultP2P{}

	//get all assets from huobi and maps with ID coin for this fiat

	//get all assets from binance for this fiat
	currencyarr := getdatahuobi.GetCurrencyHuobi(fiat)
	//log.Println(currencyarr)

	//get pair for rate

	pairmap := getdatahuobi.GetPairFromJSONHuobi(fiat)

	var wg sync.WaitGroup
	for _, a := range currencyarr {
		wg.Add(1)
		go func(a string) {
			defer wg.Done()
			arr_val := getResultP2P3TMBHH(a, fiat, pairmap, huobi, binance)
			allOrders = append(allOrders, arr_val)

		}(a)
	}
	wg.Wait()

	for _, j := range allOrders {
		for _, i := range j {
			if (i.TotalAdvBuy > 0) && (i.TotalAdvSell > 0) {
				log.Printf("%+v\n", i)
				//result.SaveResultJsonFile(fiat, i, "3steps_tt")
				//log.Printf("3 steps taker taker. Fiat - %s, Result - %v", fiat, i)
				if (i.Profit) && (i.ProfitPercet >= huobi.PercentUser) {
					result.FormatMessageAndSend(i, "You are Taker", "You areTaker")
				}
			}
		}
	}
}

func getResultP2P3TMBHH(a, fiat string, pair map[string][]string,
	huobi getinfohuobi.ParametersHuobi, binance workingbinance.ParametersBinance) []result.ResultP2P {
	//coinidmap := workinghuobi.GetCoinIDHuobo(fiat)
	var resultP2PArr []result.ResultP2P
	pair_assets := pair[strings.ToLower(a)]
	//order_buy := getdatahuobi.GetDataP2PHuobi(coinidmap[a], coinidmap[fiat], "sell", huobi)
	order_buy := getdata.GetDataP2PBinance(a, fiat, "Buy", binance)
	if len(order_buy.Data) > 0 {
		var transAmountFloat float64
		if binance.TransAmount != "" {
			tmpTransAmountFloat, err := strconv.ParseFloat(binance.TransAmount, 64)
			if err != nil {
				log.Println("Can't convert transAmount", err)
			}
			transAmountFloat = tmpTransAmountFloat
		} else {
			//log.Println(a, " - ", fiat, "order_buy.Data[0].MaxTradeLimit", order_buy.Data[0].MaxTradeLimit)
			tmpTransAmountFloat, err := strconv.ParseFloat(order_buy.Data[0].Adv.DynamicMaxSingleTransAmount, 64)
			if err != nil {
				log.Println("Can't convert MaxTradeLimit", err)
			}
			transAmountFloat = tmpTransAmountFloat
			huobi.Amount = strconv.Itoa(int(transAmountFloat))
			log.Println("New transAmount because didn't enter amount in beginer", huobi.Amount)
		}

		price_b := order_buy.Data[0].Adv.Price
		//fmt.Printf("%+v\n", order_buy)

		transAmountFirst := transAmountFloat / price_b
		//second step

		pair_rate := getdatahuobi.GetRatePairHuobi(pair_assets)

		var wg sync.WaitGroup
		for p := range pair_rate {
			wg.Add(1)

			go func(p string) {
				defer wg.Done()
				value := printResultP2P3TMBHH(p, a, fiat, transAmountFirst, price_b,
					pair_rate, order_buy, huobi, binance)
				resultP2PArr = append(resultP2PArr, value)
			}(p)

		}
		wg.Wait()
		return resultP2PArr
	} else {
		log.Printf("Can't find on Huobi for buy(Taker - Taker) %s - %s", fiat, a)
		return resultP2PArr
	}

}

func printResultP2P3TMBHH(p, a, fiat string, transAmountFirst, price_b float64, pair_rate map[string]float64,
	order_buy getinfobinance.AdvertiserAdv, huobi getinfohuobi.ParametersHuobi,
	binance workingbinance.ParametersBinance) result.ResultP2P {

	coinidmap := workinghuobi.GetCoinIDHuobo(fiat)

	profitResult := result.ResultP2P{}
	var transAmountSecond float64
	var assetSell string
	if strings.HasPrefix(p, strings.ToLower(a)) {
		transAmountSecond = (transAmountFirst * pair_rate[p])
		assetSell = p[len(a):]
	} else {
		transAmountSecond = (transAmountFirst / pair_rate[p])
		assetSell = p[:(len(p) - len(a))]
	}
	//third steps
	order_sell := getdatahuobi.GetDataP2PHuobi(coinidmap[strings.ToUpper(assetSell)], coinidmap[fiat],
		"sell", huobi)
	//fmt.Println("FIAT", fiat, "COIN", a, "\n")
	//fmt.Printf("%+v\n", order_sell)
	if len(order_sell.Data) == 0 {
		log.Printf("Can't find for sell on Huobi (Taker-Taker) %s - %s", fiat, a)
		return profitResult
	}
	price_s, _ := strconv.ParseFloat(order_sell.Data[0].Price, 64)

	transAmountThird := price_s * transAmountSecond

	transAmountFloat, err := strconv.ParseFloat(huobi.Amount, 64)
	if err != nil {
		log.Printf("Problem with convert transAmount to float, err - %v", err)
	}
	profitResult.Market.First = "Binance"
	profitResult.Merchant.FirstMerch = (binance.PublisherType == "merchant")
	profitResult.Market.Second = "Huobi"
	profitResult.Market.Third = "Huobi"
	profitResult.Merchant.ThirdMerch = (huobi.IsMerchant == "true")
	profitResult.Profit = transAmountThird > transAmountFloat
	profitResult.DataTime = time.Now()
	profitResult.Fiat = fiat
	profitResult.AssetsBuy = a
	profitResult.PriceAssetsBuy = price_b
	profitResult.PaymentBuy = result.PaymentMetods(order_buy)
	profitResult.LinkAssetsBuy = fmt.Sprintf("https://p2p.binance.com/en/trade/all-payments/%v?fiat=%v", a, fiat)
	profitResult.Pair = p
	profitResult.PricePair = pair_rate[p]
	profitResult.LinkMarket = result.ReturnLinkMarketHuobi(strings.ToLower(a), strings.ToLower(p))
	profitResult.AssetsSell = assetSell
	profitResult.PriceAssetsSell = price_s
	profitResult.PaymentSell = result.PaymentMetodsHuobi(order_sell)
	profitResult.LinkAssetsSell = fmt.Sprintf("https://www.huobi.com/en-us/fiat-crypto/trader/%s", strconv.Itoa(order_sell.Data[0].UID))
	profitResult.ProfitValue = transAmountThird - transAmountFloat
	profitResult.ProfitPercet = (((transAmountThird - transAmountFloat) / transAmountFloat) * 100)
	profitResult.TotalAdvBuy = order_buy.Total
	profitResult.TotalAdvSell = order_sell.TotalCount
	profitResult.AdvNoBuy = order_buy.Data[0].Adv.AdvNo
	profitResult.AdvNoSell = strconv.Itoa(order_sell.Data[0].UID)
	//fmt.Printf("%s - %s %+v\n", fiat, a, profitResult)
	return profitResult
}
