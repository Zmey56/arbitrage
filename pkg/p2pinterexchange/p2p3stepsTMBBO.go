package p2pinterexchange

import (
	"fmt"
	"github.com/Zmey56/arbitrage/pkg/getdata"
	"github.com/Zmey56/arbitrage/pkg/getdatahuobi"
	"github.com/Zmey56/arbitrage/pkg/getdataokx"
	"github.com/Zmey56/arbitrage/pkg/getinfobinance"
	"github.com/Zmey56/arbitrage/pkg/getinfookx"
	"github.com/Zmey56/arbitrage/pkg/result"
	"github.com/Zmey56/arbitrage/pkg/workingbinance"
	"log"
	"strconv"
	"strings"
	"sync"
	"time"
)

func P2P3stepsTMBBO(fiat string, binance workingbinance.ParametersBinance, okx getinfookx.ParametersOKX) {

	//get all assets and pair from binance for this fiat

	//get pair for rate

	pair := getinfobinance.GetPairFromJSON(fiat)

	for a, p := range pair {
		log.Println("A", a, "PAIR", p)
	}

	//get information about orders with binance
	var wg sync.WaitGroup
	for a, _ := range pair {
		wg.Add(1)
		go func(a string) {
			defer wg.Done()
			log.Println("P", pair[a], "A", a)
			getResultP2P3TMBBO(a, fiat, pair[a], binance, okx)
		}(a)
	}
	wg.Wait()

}

func getResultP2P3TMBBO(a, fiat string, pair []string, binance workingbinance.ParametersBinance,
	okx getinfookx.ParametersOKX) {

	//get all assets from binance for this fiat
	currencyarr := getdatahuobi.GetCurrencyHuobi(fiat)
	pair_assets := CheckMatchesPair(a, pair, currencyarr)
	log.Println("pair_assets", pair_assets)

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
	}
	if len(order_buy.Data) > 0 && len(pair_assets) > 0 {
		log.Println("OK")
		price_b := order_buy.Data[0].Adv.Price

		transAmountFirst := transAmountFloat / price_b
		//second step

		pair_rate := getinfobinance.GetRatePair(pair_assets)
		log.Println("pair_rate", pair_rate)

		var wg sync.WaitGroup

		for p := range pair_rate {
			wg.Add(1)

			go func(p string) {
				defer wg.Done()

				printResultP2P3TMBBO(p, a, fiat, transAmountFirst, price_b,
					pair_rate, order_buy, binance, okx)
			}(p)

		}
		wg.Wait()
	} else {
		log.Printf("Order buy is empty, fiat - %s, assets - %s, param %+v\n", fiat, a, binance)
	}

}

func printResultP2P3TMBBO(p, a, fiat string, transAmountFirst, price_b float64, pair_rate map[string]float64,
	order_buy getinfobinance.AdvertiserAdv, binance workingbinance.ParametersBinance,
	okx getinfookx.ParametersOKX) {

	//coinidmap := workinghuobi.GetCoinIDHuobo(fiat)

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
		okx.Amount = binance.TransAmount
		transAmountFirst, _ = strconv.ParseFloat(binance.TransAmount, 64)
	} else {
		transAmountFirst, _ = strconv.ParseFloat(order_buy.Data[0].Adv.DynamicMaxSingleTransAmount, 64)
		okx.Amount = order_buy.Data[0].Adv.DynamicMaxSingleTransAmount
	}
	order_sell := getdataokx.GetDataP2POKXBuy(strings.ToLower(assetSell), fiat, okx)

	if len(order_sell.Data.Buy) > 0 {

		price_s, _ := strconv.ParseFloat(order_sell.Data.Sell[0].Price, 64)
		transAmountThird := price_s * transAmountSecond

		profitResult.Amount = binance.TransAmount
		profitResult.Market.First = "Binance"
		profitResult.Merchant.FirstMerch = (binance.PublisherType == "merchant")
		profitResult.User.FirstUser = "Taker"
		profitResult.Market.Second = "Binance"
		profitResult.Market.Third = "OKX"
		profitResult.Merchant.ThirdMerch = (okx.IsMerchant == "true")
		profitResult.User.ThirdUser = "Maker"
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
		profitResult.PaymentSell = order_sell.Data.Sell[0].PaymentMethods
		profitResult.LinkAssetsSell = fmt.Sprintf("https://www.okx.com/p2p-markets/%s/buy-%s/", strings.ToLower(assetSell), strings.ToLower(a))
		profitResult.ProfitValue = transAmountThird - transAmountFirst
		profitResult.ProfitPercet = (((transAmountThird - transAmountFirst) / transAmountFirst) * 100)
		profitResult.TotalAdvBuy = order_buy.Total
		profitResult.TotalAdvSell = order_sell.Data.Total
		//log.Printf("TM %+v\n", profitResult)
		//return profitResult
		result.CheckResultSaveSend(profitResult.User.FirstUser, profitResult.User.ThirdUser, binance.Border, binance.PercentUser, profitResult)
	} else {
		log.Printf("Order sell is empty, fiat - %s, assets - %s, param %+v\n", fiat, assetSell, okx)
	}

}
