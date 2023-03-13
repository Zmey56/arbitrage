package p2phuobi

import (
	"fmt"
	"github.com/Zmey56/arbitrage/pkg/getdatahuobi"
	"github.com/Zmey56/arbitrage/pkg/getinfohuobi"
	"github.com/Zmey56/arbitrage/pkg/result"
	"github.com/Zmey56/arbitrage/pkg/workinghuobi"
	"log"
	"strconv"
	"strings"
	"sync"
	"time"
)

func P2P3stepsTakerMakerHuobi(fiat string, paramUser getinfohuobi.ParametersHuobi) {

	//get all assets from huobi and maps with ID coin for this fiat

	//get all assets from binance for this fiat
	currencyarr := getdatahuobi.GetCurrencyHuobi(fiat)

	//get pair for rate

	pairmap := getdatahuobi.GetPairFromJSONHuobi(fiat)

	var wg sync.WaitGroup
	for _, a := range currencyarr {
		wg.Add(1)
		go func(a string) {
			defer wg.Done()
			getResultP2P3TM(a, fiat, pairmap, paramUser)

		}(a)
	}
	wg.Wait()

}

func getResultP2P3TM(a, fiat string, pair map[string][]string,
	paramUser getinfohuobi.ParametersHuobi) {
	coinidmap := workinghuobi.GetCoinIDHuobo(fiat)
	pair_assets := pair[strings.ToLower(a)]
	order_buy := getdatahuobi.GetDataP2PHuobi(coinidmap[a], coinidmap[fiat], "sell", paramUser)

	if len(order_buy.Data) > 0 {
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
				log.Println("Can't convert dynamicMaxSingleTransAmount", err)
			}
			transAmountFloat = tmpTransAmountFloat
			paramUser.Amount = strconv.Itoa(int(transAmountFloat))
		}

		price_b, _ := strconv.ParseFloat(order_buy.Data[0].Price, 64)

		transAmountFirst := transAmountFloat / price_b
		//second step

		pair_rate := getdatahuobi.GetRatePairHuobi(pair_assets)

		var wg sync.WaitGroup
		for p := range pair_rate {
			wg.Add(1)

			go func(p string) {
				defer wg.Done()
				printResultP2P3TM(p, a, fiat, transAmountFirst, price_b,
					pair_rate, order_buy, paramUser)
			}(p)

		}
		wg.Wait()
	} else {
		log.Printf("Order buy is empty, fiat - %s, assets - %s, param %+v\n", fiat, a, paramUser)
	}

}

func printResultP2P3TM(p, a, fiat string, transAmountFirst, price_b float64, pair_rate map[string]float64,
	order_buy getdatahuobi.Huobi, paramUser getinfohuobi.ParametersHuobi) {

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
		"sell", paramUser)

	if len(order_sell.Data) == 0 {
		log.Printf("Order sell is empty, fiat - %s, assets - %s, param %+v\n", fiat, a, paramUser)
	} else {
		price_s, _ := strconv.ParseFloat(order_sell.Data[0].Price, 64)

		transAmountThird := price_s * transAmountSecond

		transAmountFloat, err := strconv.ParseFloat(paramUser.Amount, 64)
		if err != nil {
			log.Printf("Problem with convert transAmount to float, err - %v", err)
		}
		profitResult.Amount = paramUser.Amount
		profitResult.Market.First = "Huobi"
		profitResult.Merchant.FirstMerch = (paramUser.IsMerchant == "true")
		profitResult.User.FirstUser = "Taker"
		profitResult.Market.Second = "Huobi"
		profitResult.Market.Third = "Huobi"
		profitResult.User.ThirdUser = "Maker"
		profitResult.Merchant.ThirdMerch = (paramUser.IsMerchant == "true")
		profitResult.Profit = transAmountThird > transAmountFloat
		profitResult.DataTime = time.Now()
		profitResult.Fiat = fiat
		profitResult.AssetsBuy = a
		profitResult.PriceAssetsBuy = price_b
		profitResult.PaymentBuy = result.PaymentMetodsHuobi(order_buy)
		profitResult.LinkAssetsBuy = fmt.Sprintf("https://www.huobi.com/en-us/fiat-crypto/trade/buy-%s-%s/", strings.ToLower(a), strings.ToLower(fiat))
		profitResult.Pair = strings.ToUpper(p)
		profitResult.PricePair = pair_rate[p]
		profitResult.LinkMarket = result.ReturnLinkMarketHuobi(strings.ToLower(a), strings.ToLower(p))
		profitResult.AssetsSell = strings.ToUpper(assetSell)
		profitResult.PriceAssetsSell = price_s
		profitResult.PaymentSell = result.PaymentMetodsHuobi(order_sell)
		profitResult.LinkAssetsSell = fmt.Sprintf("https://www.huobi.com/en-us/fiat-crypto/trade/buy-%s-%s/", assetSell, strings.ToLower(fiat))
		profitResult.ProfitValue = transAmountThird - transAmountFloat
		profitResult.ProfitPercet = (((transAmountThird - transAmountFloat) / transAmountFloat) * 100)
		profitResult.TotalAdvBuy = order_buy.TotalCount
		profitResult.TotalAdvSell = order_sell.TotalCount
		profitResult.AdvNoBuy = strconv.Itoa(order_buy.Data[0].UID)
		profitResult.AdvNoSell = strconv.Itoa(order_sell.Data[0].UID)
		//return profitResult
		result.CheckResultSaveSend(profitResult.User.FirstUser, profitResult.User.ThirdUser, paramUser.Border, paramUser.PercentUser, profitResult)
	}
}
