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

func P2P3stepsTakerTakerHuobi(fiat string, paramUser getinfohuobi.ParametersHuobi) {
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
			getResultP2P3TT(a, fiat, pairmap, paramUser)

		}(a)
	}
	wg.Wait()

}

func getResultP2P3TT(a, fiat string, pair map[string][]string,
	paramUser getinfohuobi.ParametersHuobi) {
	coinidmap := workinghuobi.GetCoinIDHuobo(fiat)
	pair_assets := pair[strings.ToLower(a)]

	//log.Printf("coinidmap %+v\n", coinidmap)
	//log.Println("coinidmap[a]", a, "-", coinidmap[a], "fiat", fiat, " - ", coinidmap[fiat])
	//log.Println("coinidmap[a], coinidmap[fiat], paramUser", coinidmap[a], coinidmap[fiat], "sell", paramUser)
	order_buy := getdatahuobi.GetDataP2PHuobi(coinidmap[a], coinidmap[fiat], "sell", paramUser)
	//log.Printf("%+v\n", order_buy)

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
				log.Println("Can't convert MaxTradeLimit", err)
			}
			transAmountFloat = tmpTransAmountFloat
			paramUser.Amount = strconv.Itoa(int(transAmountFloat))
		}

		price_b, _ := strconv.ParseFloat(order_buy.Data[0].Price, 64)

		transAmountFirst := transAmountFloat / price_b
		//second step

		//log.Println("pair_assets", pair_assets)
		pair_rate := getdatahuobi.GetRatePairHuobi(pair_assets)
		//log.Println("pair_rate", pair_rate)

		var wg sync.WaitGroup
		for p := range pair_rate {
			wg.Add(1)

			go func(p string) {
				defer wg.Done()
				printResultP2P3TT(p, a, fiat, transAmountFirst, price_b,
					pair_rate, order_buy, paramUser)
			}(p)

		}
		wg.Wait()
	} else {
		log.Printf("Order buy is empty, fiat - %s, assets - %s, param %+v\n", fiat, a, paramUser)
	}

}

func printResultP2P3TT(p, a, fiat string, transAmountFirst, price_b float64, pair_rate map[string]float64,
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
	//log.Println("coinidmap[strings.ToUpper(assetSell)], coinidmap[fiat], paramUser", coinidmap[strings.ToUpper(assetSell)], coinidmap[fiat],
	//	"buy", paramUser)
	order_sell := getdatahuobi.GetDataP2PHuobi(coinidmap[strings.ToUpper(assetSell)], coinidmap[fiat],
		"buy", paramUser)
	//log.Printf("%+v\n\n", order_sell)

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
		profitResult.Merchant.ThirdMerch = (paramUser.IsMerchant == "true")
		profitResult.User.ThirdUser = "Taker"
		profitResult.Profit = transAmountThird > transAmountFloat
		profitResult.DataTime = time.Now()
		profitResult.Fiat = fiat
		profitResult.AssetsBuy = a
		profitResult.PriceAssetsBuy = price_b
		profitResult.PaymentBuy = result.PaymentMetodsHuobi(order_buy)
		profitResult.LinkAssetsBuy = fmt.Sprintf("https://www.huobi.com/en-us/fiat-crypto/trade/buy-%s-%s/", strings.ToLower(a), strings.ToLower(fiat))
		profitResult.Pair = p
		profitResult.PricePair = pair_rate[p]
		profitResult.LinkMarket = result.ReturnLinkMarketHuobi(strings.ToLower(a), strings.ToLower(p))
		profitResult.AssetsSell = assetSell
		profitResult.PriceAssetsSell = price_s
		profitResult.PaymentSell = result.PaymentMetodsHuobi(order_sell)
		profitResult.LinkAssetsSell = fmt.Sprintf("https://p2p.binance.com/en/trade/sell/%v?fiat=%v", assetSell, fiat)
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
