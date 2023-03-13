package p2pinterexchange

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

	//get all assets and pair from binance for this fiat

	//get pair for rate

	pair := getinfobinance.GetPairFromJSON(fiat)

	//get information about orders with binance
	var wg sync.WaitGroup
	for a, p := range pair {
		wg.Add(1)
		go func(a string) {
			defer wg.Done()
			getResultP2P3TMBBH(a, fiat, p, binance, huobi)
		}(a)
	}
	wg.Wait()

}

func getResultP2P3TMBBH(a, fiat string, pair []string, binance workingbinance.ParametersBinance,
	huobi getinfohuobi.ParametersHuobi) {

	//get all assets from binance for this fiat
	currencyarr := getdatahuobi.GetCurrencyHuobi(fiat)
	pair_assets := CheckMatchesPair(a, pair, currencyarr)

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

		price_b := order_buy.Data[0].Adv.Price

		transAmountFirst := transAmountFloat / price_b
		//second step

		pair_rate := getinfobinance.GetRatePair(pair_assets)

		var wg sync.WaitGroup

		for p := range pair_rate {
			wg.Add(1)

			go func(p string) {
				defer wg.Done()

				printResultP2P3TMBBH(p, a, fiat, transAmountFirst, price_b,
					pair_rate, order_buy, binance, huobi)
			}(p)

		}
		wg.Wait()
	} else {
		log.Printf("Order buy is empty, fiat - %s, assets - %s, param %+v\n", fiat, a, binance)
	}

}

func printResultP2P3TMBBH(p, a, fiat string, transAmountFirst, price_b float64, pair_rate map[string]float64,
	order_buy getinfobinance.AdvertiserAdv, binance workingbinance.ParametersBinance,
	huobi getinfohuobi.ParametersHuobi) {

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
	order_sell := getdatahuobi.GetDataP2PHuobi(coinidmap[strings.ToUpper(assetSell)], coinidmap[fiat],
		"sell", huobi)
	if len(order_sell.Data) > 0 {

		price_s, _ := strconv.ParseFloat(order_sell.Data[0].Price, 64)
		transAmountThird := price_s * transAmountSecond

		profitResult.Amount = binance.TransAmount
		profitResult.Market.First = "Binance"
		profitResult.Merchant.FirstMerch = (binance.PublisherType == "merchant")
		profitResult.User.FirstUser = "Taker"
		profitResult.Market.Second = "Binance"
		profitResult.Market.Third = "Huobi"
		profitResult.Merchant.ThirdMerch = (huobi.IsMerchant == "true")
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
		profitResult.PaymentSell = result.PaymentMetodsHuobi(order_sell)
		profitResult.LinkAssetsSell = fmt.Sprintf("https://www.huobi.com/en-us/fiat-crypto/trader/%s", strconv.Itoa(order_sell.Data[0].UID))
		profitResult.ProfitValue = transAmountThird - transAmountFirst
		profitResult.ProfitPercet = (((transAmountThird - transAmountFirst) / transAmountFirst) * 100)
		profitResult.TotalAdvBuy = order_buy.Total
		profitResult.TotalAdvSell = order_sell.TotalCount
		//log.Printf("TM %+v\n", profitResult)
		//return profitResult
		result.CheckResultSaveSend(profitResult.User.FirstUser, profitResult.User.ThirdUser, binance.Border, binance.PercentUser, profitResult)
	} else {
		log.Printf("Order sell is empty, fiat - %s, assets - %s, param %+v\n", fiat, assetSell, huobi)
	}

}
