package p2pinterexchange

import (
	"fmt"
	"github.com/Zmey56/arbitrage/pkg/getdata"
	"github.com/Zmey56/arbitrage/pkg/getdatabh"
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

func P2P3stepsTMHBB(fiat string, binance workingbinance.ParametersBinance,
	huobi getinfohuobi.ParametersHuobi) {

	//get pair for rate

	pair := getdatabh.GetPairFromJSONHBB(fiat)

	//get information about orders with binance
	var wg sync.WaitGroup
	for a, p := range pair {
		wg.Add(1)
		go func(a string) {
			defer wg.Done()
			getResultP2P3TMHBB(a, fiat, p, binance, huobi)
		}(a)
	}
	wg.Wait()

}

func getResultP2P3TMHBB(a, fiat string, pair []string, binance workingbinance.ParametersBinance,
	huobi getinfohuobi.ParametersHuobi) {

	coinidmap := workinghuobi.GetCoinIDHuobo(fiat)
	//first step
	order_buy := getdatahuobi.GetDataP2PHuobi(coinidmap[a], coinidmap[fiat], "sell", huobi)
	var transAmountFloat float64
	if huobi.Amount != "" {
		tmpTransAmountFloat, err := strconv.ParseFloat(huobi.Amount, 64)
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
		binance.TransAmount = strconv.Itoa(int(transAmountFloat))
		log.Println("New transAmount because didn't enter amount in beginer", binance.TransAmount)
	}

	if len(order_buy.Data) > 0 {

		price_b, _ := strconv.ParseFloat(order_buy.Data[0].Price, 64)

		transAmountFirst := transAmountFloat / price_b
		//second step

		pair_rate := getinfobinance.GetRatePair(pair)

		var wg sync.WaitGroup
		for p := range pair_rate {
			wg.Add(1)

			go func(p string) {
				defer wg.Done()
				printResultP2P3TMHBB(p, a, fiat, transAmountFirst, price_b,
					pair_rate, order_buy, binance, huobi)
			}(p)

		}
		wg.Wait()
	} else {
		log.Printf("Order buy is empty, fiat - %s, assets - %s, param %+v\n", fiat, a, huobi)
	}

}

func printResultP2P3TMHBB(p, a, fiat string, transAmountFirst, price_b float64, pair_rate map[string]float64,
	order_buy getdatahuobi.Huobi, binance workingbinance.ParametersBinance,
	huobi getinfohuobi.ParametersHuobi) {

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
		"Buy", binance)
	if len(order_sell.Data) == 0 {
		log.Printf("Order sell is empty, fiat - %s, assets - %s, param %+v\n", fiat, assetSell, binance)

	} else {
		price_s := order_sell.Data[0].Adv.Price

		transAmountThird := price_s * transAmountSecond

		transAmountFloat, err := strconv.ParseFloat(binance.TransAmount, 64)
		if err != nil {
			log.Printf("Problem with convert transAmount to float, err - %v", err)
		}
		profitResult.Amount = huobi.Amount
		profitResult.Market.First = "Huobi"
		profitResult.Merchant.FirstMerch = (huobi.IsMerchant == "true")
		profitResult.User.FirstUser = "Taker"
		profitResult.Market.Second = "Binance"
		profitResult.Market.Third = "Binance"
		profitResult.Merchant.ThirdMerch = (binance.PublisherType == "merchant")
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
		profitResult.LinkMarket = result.ReturnLinkMarket(a, p)
		profitResult.AssetsSell = assetSell
		profitResult.PriceAssetsSell = price_s
		profitResult.PaymentSell = result.PaymentMetods(order_sell)
		profitResult.LinkAssetsSell = fmt.Sprintf("https://p2p.binance.com/en/trade/all-payments/%v?fiat=%v",
			assetSell, fiat)
		profitResult.ProfitValue = transAmountThird - transAmountFloat
		profitResult.ProfitPercet = (((transAmountThird - transAmountFloat) / transAmountFloat) * 100)
		profitResult.TotalAdvBuy = order_buy.TotalCount
		profitResult.TotalAdvSell = order_sell.Total
		profitResult.AdvNoBuy = strconv.Itoa(order_buy.Data[0].UID)
		profitResult.AdvNoSell = order_sell.Data[0].Adv.AdvNo
		//return profitResult
		result.CheckResultSaveSend(profitResult.User.FirstUser, profitResult.User.ThirdUser, binance.Border, binance.PercentUser, profitResult)
	}
}
