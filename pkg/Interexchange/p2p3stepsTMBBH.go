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

func P2P3stepsTMBBH(fiat string, binance workingbinance.ParametersBinance, huobi getinfohuobi.ParametersHuobi) {
	log.Println(binance)
	allOrders := [][]result.ResultP2P{}
	//get all assets from binance for this fiat

	assets := getdata.GetAssets(fiat)
	assets_symbol := make([]string, 0, len(assets))
	assets_name := make([]string, 0, len(assets))

	for k, v := range assets {
		assets_symbol = append(assets_symbol, k)
		assets_name = append(assets_name, v)
	}

	//get pair for rate

	pair := getinfobinance.GetPairFromJSON(fiat)

	//get information about orders with binance
	var wg sync.WaitGroup
	for _, a := range assets_symbol {
		wg.Add(1)
		go func(a string) {
			defer wg.Done()
			arr_val := getResultP2P3TMBBH(a, fiat, pair, binance, huobi)
			allOrders = append(allOrders, arr_val)
		}(a)
	}
	wg.Wait()

	for _, j := range allOrders {
		for _, i := range j {
			if (i.TotalAdvBuy > 0) && (i.TotalAdvSell > 0) {
				//result.SaveResultJsonFile(fiat, i, "3steps_tt")
				if (i.Profit) && (i.ProfitPercet >= binance.PercentUser) {
					result.FormatMessageAndSend(i, "You are Taker", "You areTaker")
				}
			}
		}
	}
}

func getResultP2P3TMBBH(a, fiat string, pair map[string][]string,
	binance workingbinance.ParametersBinance, huobi getinfohuobi.ParametersHuobi) []result.ResultP2P {

	var resultP2PArr []result.ResultP2P
	pair_ass := pair[a]
	//get all assets from binance for this fiat
	currencyarr := getdatahuobi.GetCurrencyHuobi(fiat)
	pair_assets := CheckMatchesPair(a, pair_ass, currencyarr)

	//first step
	order_buy := getdata.GetDataP2PBinance(a, fiat, "Buy", binance)
	var transAmountFloat float64
	if binance.TransAmount != "" {
		tmpTransAmountFloat, err := strconv.ParseFloat(binance.TransAmount, 64)
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
		binance.TransAmount = strconv.Itoa(int(transAmountFloat))
		log.Println("New transAmount because didn't enter amount in beginer", binance.TransAmount)
	}
	if len(order_buy.Data) > 0 || len(pair_assets) > 0 {

		price_b := order_buy.Data[0].Adv.Price

		transAmountFirst := transAmountFloat / price_b
		//second step

		pair_rate := getinfobinance.GetRatePair(pair_assets)

		var wg sync.WaitGroup

		for p := range pair_rate {
			wg.Add(1)

			go func(p string) {
				defer wg.Done()

				value := printResultP2P3TMBBH(p, a, fiat, transAmountFirst, price_b,
					pair_rate, order_buy, binance, huobi)
				resultP2PArr = append(resultP2PArr, value)
			}(p)

		}
		wg.Wait()
		return resultP2PArr
	} else {
		return resultP2PArr
	}

}

func printResultP2P3TMBBH(p, a, fiat string, transAmountFirst, price_b float64, pair_rate map[string]float64,
	order_buy getinfobinance.AdvertiserAdv, binance workingbinance.ParametersBinance,
	huobi getinfohuobi.ParametersHuobi) result.ResultP2P {

	coinidmap := workinghuobi.GetCoinIDHuobo(fiat)

	profitResult := result.ResultP2P{}

	//with pair on Binance
	var transAmountSecond float64
	var assetSell string
	if strings.HasPrefix(p, a) {
		transAmountSecond = (transAmountFirst * pair_rate[p])
		assetSell = p[len(a):]
	} else {
		transAmountSecond = (transAmountFirst / pair_rate[p])
		assetSell = p[:(len(p) - len(a))]
	}

	//third staps with huobi
	if binance.TransAmount != "" {
		huobi.Amount = binance.TransAmount
		transAmountFirst, _ = strconv.ParseFloat(binance.TransAmount, 64)
	} else {
		transAmountFirst, _ = strconv.ParseFloat(order_buy.Data[0].Adv.DynamicMaxSingleTransAmount, 64)
		huobi.Amount = order_buy.Data[0].Adv.DynamicMaxSingleTransAmount
	}
	log.Println("Ass", a, "PAIR", p)
	order_sell := getdatahuobi.GetDataP2PHuobi(coinidmap[strings.ToUpper(assetSell)], coinidmap[fiat],
		"sell", huobi)
	log.Println("order_sell", order_sell.Data[0].Price, "Pair", p, "assetSell", assetSell)
	if len(order_sell.Data) == 0 {
		return profitResult
	}
	price_s, _ := strconv.ParseFloat(order_sell.Data[0].Price, 64)
	transAmountThird := price_s * transAmountSecond

	profitResult.Market.First = "Binance"
	profitResult.Merchant.FirstMerch = (binance.PublisherType == "merchant")
	profitResult.Market.Second = "Binance"
	profitResult.Market.Third = "Huobi"
	profitResult.Merchant.ThirdMerch = (huobi.IsMerchant == "true")
	profitResult.Profit = transAmountThird > transAmountFirst
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
	profitResult.PaymentSell = result.PaymentMetodsHuobi(order_sell)
	profitResult.LinkAssetsSell = fmt.Sprintf("https://www.huobi.com/en-us/fiat-crypto/trader/%s", strconv.Itoa(order_sell.Data[0].UID))
	profitResult.ProfitValue = transAmountThird - transAmountFirst
	profitResult.ProfitPercet = (((transAmountThird - transAmountFirst) / transAmountFirst) * 100)
	profitResult.TotalAdvBuy = order_buy.Total
	profitResult.TotalAdvSell = order_sell.TotalCount
	return profitResult
}
