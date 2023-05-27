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

func P2P3stepsTakerTakerOKX(fiat string, paramUser getinfookx.ParametersOKX) {

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
			getResultP2P3TT(a, fiat, pairmap, paramUser)

		}(a)
	}
	wg.Wait()

}

func getResultP2P3TT(a, fiat string, pair map[string][]string,
	paramUser getinfookx.ParametersOKX) {
	order_buy := getdataokx.GetDataP2POKXBuy(fiat, a, paramUser)

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

		pair_rate := getdataokx.GetRatePairOKX(pair[a])

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
	}

}

func printResultP2P3TT(p, a, fiat string, transAmountFirst, price_b float64, pair_rate map[string]float64,
	order_buy getdataokx.OKXBuy, paramUser getinfookx.ParametersOKX) {

	profitResult := result.ResultP2P{}
	var transAmountSecond float64
	var assetSell string
	splitPair := strings.Split(p, "-")
	if a == splitPair[0] {
		transAmountSecond = (transAmountFirst * pair_rate[p])
		assetSell = splitPair[1]
	} else {
		transAmountSecond = (transAmountFirst / pair_rate[p])
		assetSell = splitPair[0]
	}
	//third steps

	order_sell := getdataokx.GetDataP2POKXSell(strings.ToLower(fiat), strings.ToLower(assetSell), paramUser)

	if len(order_sell.Data.Buy) > 0 {
		price_s, _ := strconv.ParseFloat(order_sell.Data.Buy[0].Price, 64)

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
		profitResult.User.ThirdUser = "Taker"
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
		profitResult.PaymentSell = order_sell.Data.Buy[0].PaymentMethods
		profitResult.LinkAssetsSell = fmt.Sprintf("https://www.okx.com/p2p-markets/%s/sell-%s", fiat, assetSell)
		profitResult.ProfitValue = transAmountThird - transAmountFloat
		profitResult.ProfitPercet = (((transAmountThird - transAmountFloat) / transAmountFloat) * 100)
		profitResult.TotalAdvBuy = order_buy.Data.Total
		profitResult.TotalAdvSell = order_sell.Data.Total
		profitResult.AdvNoBuy = order_buy.Data.Sell[0].ID
		profitResult.AdvNoSell = order_sell.Data.Buy[0].ID
		result.CheckResultSaveSend(profitResult.User.FirstUser, profitResult.User.ThirdUser, paramUser.Border, paramUser.PercentUser, profitResult)
	}
}
