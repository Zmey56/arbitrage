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

func P2P3stepsTakerTaker(fiat string, paramUser workingbinance.ParametersBinance) {
	//get all assets and pair from binance for this fiat

	pair := getinfobinance.GetPairFromJSON(fiat)

	var wg sync.WaitGroup
	for a, _ := range pair {
		wg.Add(1)
		go func(a string) {
			defer wg.Done()
			GetResultP2P3TT(a, fiat, pair, paramUser)
		}(a)
	}
	wg.Wait()

}

func GetResultP2P3TT(a, fiat string, pair map[string][]string, paramUser workingbinance.ParametersBinance) {
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
				PrintResultP2P3TT(p, a, fiat, transAmountFirst, price_b,
					pair_rate, order_buy, paramUser)
			}(p)

		}
		wg.Wait()
	} else {
		log.Printf("Order buy is empty, fiat - %s, assets - %s, param %+v\n", fiat, a, paramUser)
	}

}

func PrintResultP2P3TT(p, a, fiat string, transAmountFirst, price_b float64, pair_rate map[string]float64,
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
		"Sell", paramUser)
	if len(order_sell.Data) == 0 {
		log.Printf("Order sell is empty, fiat - %s, assets - %s, param %+v\n", fiat, a, paramUser)
	} else {
		price_s := order_sell.Data[0].Adv.Price

		transAmountThird := price_s * transAmountSecond

		transAmountFloat, err := strconv.ParseFloat(paramUser.TransAmount, 64)
		if err != nil {
			log.Printf("Problem with convert transAmount to float, err - %v", err)
		}
		profitResult.Amount = paramUser.TransAmount
		profitResult.Market.First = "Binance"
		profitResult.User.FirstUser = "Taker"
		profitResult.Merchant.FirstMerch = (paramUser.PublisherType == "merchant")
		profitResult.Market.Second = "Binance"
		profitResult.Market.Third = "Binance"
		profitResult.User.ThirdUser = "Taker"
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
		profitResult.LinkAssetsSell = fmt.Sprintf("https://p2p.binance.com/en/trade/sell/%v?fiat=%v&payment=ALL",
			assetSell, fiat)
		profitResult.ProfitValue = transAmountThird - transAmountFloat
		profitResult.ProfitPercet = (((transAmountThird - transAmountFloat) / transAmountFloat) * 100)
		profitResult.TotalAdvBuy = order_buy.Total
		profitResult.TotalAdvSell = order_sell.Total
		//return profitResult

		result.CheckResultSaveSend(profitResult.User.FirstUser, profitResult.User.ThirdUser, paramUser.Border, paramUser.PercentUser, profitResult)
	}
}
