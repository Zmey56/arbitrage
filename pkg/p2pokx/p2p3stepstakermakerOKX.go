package p2pokx

import (
	"fmt"
	"github.com/Zmey56/arbitrage/pkg/getdataokx"
	"github.com/Zmey56/arbitrage/pkg/getinfookx"
	"github.com/Zmey56/arbitrage/pkg/result"
	"log"
	"strconv"
	"strings"
	"sync"
	"time"
)

func P2P3stepsTakerMakerOKX(fiat string, paramUser getinfookx.ParametersOKX) {

	//get pair for rate
	pairmap := getdataokx.GetPairFromJSONOKX(fiat)

	//get all assets from binance for this fiat
	currencyarr := []string{}

	for i, _ := range pairmap {
		currencyarr = append(currencyarr, i)
	}

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
	paramUser getinfookx.ParametersOKX) {
	//coinidmap := workinghuobi.GetCoinIDHuobo(fiat)
	//pair_assets := pair[strings.ToLower(a)]

	//log.Printf("coinidmap %+v\n", coinidmap)
	//log.Println("COIN", a, "fiat", fiat)
	//log.Println("coinidmap[a], coinidmap[fiat], paramUser", coinidmap[a], coinidmap[fiat], "sell", paramUser)
	order_buy := getdataokx.GetDataP2POKXBuy(fiat, a, paramUser)
	//log.Printf("%+v\n", order_buy)

	if len(order_buy.Data.Sell) > 0 {
		var transAmountFloat float64
		if paramUser.Amount != "" {
			tmpTransAmountFloat, err := strconv.ParseFloat(paramUser.Amount, 64)
			if err != nil {
				log.Println("Can't convert transAmount", err)
			}
			transAmountFloat = tmpTransAmountFloat
		} else {
			tmpTransAmountFloat, err := strconv.ParseFloat(order_buy.Data.Sell[0].QuoteMaxAmountPerOrder, 64)
			if err != nil {
				log.Println("Can't convert MaxTradeLimit", err)
			}
			transAmountFloat = tmpTransAmountFloat
			paramUser.Amount = strconv.Itoa(int(transAmountFloat))
		}

		price_b, _ := strconv.ParseFloat(order_buy.Data.Sell[0].Price, 64)

		transAmountFirst := transAmountFloat / price_b
		//second step

		//log.Println("pair_assets", pair[a])
		pair_rate := getdataokx.GetRatePairOKX(pair[a])
		//log.Println("pair_rate", fiat, "-", a, "|", pair[a], "Result", pair_rate)

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
	order_buy getdataokx.OKXBuy, paramUser getinfookx.ParametersOKX) {

	//coinidmap := workinghuobi.GetCoinIDHuobo(fiat)

	profitResult := result.ResultP2P{}
	var transAmountSecond float64
	var assetSell string
	splitPair := strings.Split(p, "-")
	//log.Println("P", p, "A", a, a == splitPair[0], a == splitPair[1])
	if a == splitPair[0] {
		transAmountSecond = (transAmountFirst * pair_rate[p])
		assetSell = splitPair[1]
	} else {
		transAmountSecond = (transAmountFirst / pair_rate[p])
		assetSell = splitPair[0]
	}
	//third steps

	order_sell := getdataokx.GetDataP2POKXBuy(strings.ToLower(fiat), strings.ToLower(assetSell), paramUser)
	log.Printf("%+v\n\n", order_sell)

	if len(order_sell.Data.Sell) == 0 {
		log.Printf("Order sell is empty, fiat - %s, assets - %s, param %+v\n", fiat, a, paramUser)
	} else {
		price_s, _ := strconv.ParseFloat(order_sell.Data.Sell[0].Price, 64)

		transAmountThird := price_s * transAmountSecond

		transAmountFloat, err := strconv.ParseFloat(paramUser.Amount, 64)
		if err != nil {
			log.Printf("Problem with convert transAmount to float, err - %v", err)
		}
		profitResult.Amount = paramUser.Amount
		profitResult.Market.First = "OKX"
		profitResult.Merchant.FirstMerch = (paramUser.IsMerchant == "true")
		profitResult.User.FirstUser = "Taker"
		profitResult.Market.Second = "OKX"
		profitResult.Market.Third = "OKX"
		profitResult.Merchant.ThirdMerch = (paramUser.IsMerchant == "true")
		profitResult.User.ThirdUser = "Maker"
		profitResult.Profit = transAmountThird > transAmountFloat
		profitResult.DataTime = time.Now()
		profitResult.Fiat = fiat
		profitResult.AssetsBuy = a
		profitResult.PriceAssetsBuy = price_b
		profitResult.PaymentBuy = order_buy.Data.Sell[0].PaymentMethods
		profitResult.LinkAssetsBuy = fmt.Sprintf("https://www.okx.com/p2p-markets/%s/buy-%s/", strings.ToLower(fiat), strings.ToLower(a))
		profitResult.Pair = p
		profitResult.PricePair = pair_rate[p]
		profitResult.LinkMarket = fmt.Sprintf("https://www.okx.com/trade-spot/%s", strings.ToLower(p))
		profitResult.AssetsSell = assetSell
		profitResult.PriceAssetsSell = price_s
		profitResult.PaymentSell = order_sell.Data.Sell[0].PaymentMethods
		profitResult.LinkAssetsSell = fmt.Sprintf("https://www.okx.com/p2p-markets/%s/buy-%s/", strings.ToLower(assetSell), strings.ToLower(a))
		profitResult.ProfitValue = transAmountThird - transAmountFloat
		profitResult.ProfitPercet = (((transAmountThird - transAmountFloat) / transAmountFloat) * 100)
		profitResult.TotalAdvBuy = order_buy.Data.Total
		profitResult.TotalAdvSell = order_sell.Data.Total
		profitResult.AdvNoBuy = order_buy.Data.Sell[0].ID
		profitResult.AdvNoSell = order_sell.Data.Sell[0].ID
		//return profitResult
		result.CheckResultSaveSend(profitResult.User.FirstUser, profitResult.User.ThirdUser, paramUser.Border, paramUser.PercentUser, profitResult)
	}
}
