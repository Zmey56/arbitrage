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

func P2P3stepsTTBHH(fiat string, huobi getinfohuobi.ParametersHuobi,
	binance workingbinance.ParametersBinance) {
	//log.Println(huobi)

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
			getResultP2P3TTBHH(a, fiat, pairmap, huobi, binance)

		}(a)
	}
	wg.Wait()
}

func getResultP2P3TTBHH(a, fiat string, pair map[string][]string,
	huobi getinfohuobi.ParametersHuobi, binance workingbinance.ParametersBinance) {
	//coinidmap := workinghuobi.GetCoinIDHuobo(fiat)
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
			tmpTransAmountFloat, err := strconv.ParseFloat(order_buy.Data[0].Adv.DynamicMaxSingleTransAmount, 64)
			if err != nil {
				log.Println("Can't convert MaxTradeLimit", err)
			}
			transAmountFloat = tmpTransAmountFloat
			huobi.Amount = strconv.Itoa(int(transAmountFloat))
		}

		price_b := order_buy.Data[0].Adv.Price

		transAmountFirst := transAmountFloat / price_b
		//second step

		pair_rate := getdatahuobi.GetRatePairHuobi(pair_assets)

		var wg sync.WaitGroup
		for p := range pair_rate {
			wg.Add(1)

			go func(p string) {
				defer wg.Done()
				printResultP2P3TTBHH(p, a, fiat, transAmountFirst, price_b,
					pair_rate, order_buy, huobi, binance)
			}(p)

		}
		wg.Wait()
	} else {
		log.Printf("Order buy is empty, fiat - %s, assets - %s, param %+v\n", fiat, a, binance)
	}

}

func printResultP2P3TTBHH(p, a, fiat string, transAmountFirst, price_b float64, pair_rate map[string]float64,
	order_buy getinfobinance.AdvertiserAdv, huobi getinfohuobi.ParametersHuobi,
	binance workingbinance.ParametersBinance) {

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
		"buy", huobi)
	if len(order_sell.Data) == 0 {
		log.Printf("Order sell is empty, fiat - %s, assets - %s, param %+v\n", fiat, assetSell, huobi)
	} else {
		price_s, _ := strconv.ParseFloat(order_sell.Data[0].Price, 64)

		transAmountThird := price_s * transAmountSecond

		transAmountFloat, err := strconv.ParseFloat(huobi.Amount, 64)
		if err != nil {
			log.Printf("Problem with convert transAmount to float, err - %v", err)
		}
		profitResult.Amount = binance.TransAmount
		profitResult.Market.First = "Binance"
		profitResult.Merchant.FirstMerch = (binance.PublisherType == "merchant")
		profitResult.User.FirstUser = "Taker"
		profitResult.Market.Second = "Huobi"
		profitResult.Market.Third = "Huobi"
		profitResult.Merchant.ThirdMerch = (huobi.IsMerchant == "true")
		profitResult.User.ThirdUser = "Taker"
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
		profitResult.LinkAssetsSell = fmt.Sprintf("https://www.huobi.com/en-us/fiat-crypto/trade/buy-%s-%s/", strings.ToLower(assetSell), strings.ToLower(fiat))
		profitResult.ProfitValue = transAmountThird - transAmountFloat
		profitResult.ProfitPercet = (((transAmountThird - transAmountFloat) / transAmountFloat) * 100)
		profitResult.TotalAdvBuy = order_buy.Total
		profitResult.TotalAdvSell = order_sell.TotalCount
		profitResult.AdvNoBuy = order_buy.Data[0].Adv.AdvNo
		profitResult.AdvNoSell = strconv.Itoa(order_sell.Data[0].UID)
		//return profitResult
		result.CheckResultSaveSend(profitResult.User.FirstUser, profitResult.User.ThirdUser, binance.Border, huobi.PercentUser, profitResult)
	}
}
