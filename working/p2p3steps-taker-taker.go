package working

import (
	"fmt"
	"github.com/Zmey56/arbitrage/getdata"
	"github.com/Zmey56/arbitrage/getinfobinance"
	"github.com/Zmey56/arbitrage/interact"
	"github.com/Zmey56/arbitrage/result"
	"log"
	"strconv"
	"strings"
	"sync"
	"time"
)

func P2P3stepsTakerTaker(fiat string, paramUser interact.Parameters) {
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
			arr_val := GetResultP2P3TT(a, fiat, pair, paramUser)
			allOrders = append(allOrders, arr_val)
		}(a)
	}
	wg.Wait()

	for _, j := range allOrders {
		for _, i := range j {
			if i.Profit && i.ProfitPercet >= paramUser.PercentUser {
				log.Println("Precent", i.ProfitPercet, paramUser.PercentUser)
				result.FormatMessageAndSend(i, "Taker-Taker")
			}
		}
	}
}

func GetResultP2P3TT(a, fiat string, pair map[string][]string, paramUser interact.Parameters) []result.ResultP2P {
	//fmt.Println("====================================")
	//log.Println("ASSETS", a)
	var resultP2PArr []result.ResultP2P
	pair_assets := pair[a]
	//first step
	order_buy := getinfobinance.GetDataP2P(a, fiat, "Buy", paramUser)
	var transAmountFloat float64
	if paramUser.TransAmount != "" {
		tmpTransAmountFloat, err := strconv.ParseFloat(paramUser.TransAmount, 64)
		if err != nil {
			log.Println("Can't convert transAmount", err)
		}
		transAmountFloat = tmpTransAmountFloat
	} else {
		tmpTransAmountFloat, err := strconv.ParseFloat(order_buy.Adv.DynamicMaxSingleTransAmount, 64)
		if err != nil {
			log.Println("Can't convert dynamicMaxSingleTransAmount", err)
		}
		transAmountFloat = tmpTransAmountFloat
		paramUser.TransAmount = strconv.Itoa(int(transAmountFloat))
		log.Println("New transAmount because didn't enter amount in beginer", paramUser.TransAmount)
	}
	if order_buy.Adv.Price != 0 {

		price_b := order_buy.Adv.Price

		transAmountFirst := transAmountFloat / price_b
		//log.Printf("transAmountFirst - %v", transAmountFirst)
		//second step

		//log.Printf("Pair Assets - %s", pair_assets)
		pair_rate := getinfobinance.GetRatePair(pair_assets)

		//log.Printf("Pair Rate - %s", pair_rate)
		var wg sync.WaitGroup
		for p := range pair_rate {
			wg.Add(1)

			go func(p string) {
				defer wg.Done()
				value := PrintResultP2P3TT(p, a, fiat, transAmountFirst, price_b,
					pair_rate, order_buy, paramUser)
				resultP2PArr = append(resultP2PArr, value)
			}(p)

		}
		wg.Wait()
		return resultP2PArr
	} else {
		return resultP2PArr
	}

}

func PrintResultP2P3TT(p, a, fiat string, transAmountFirst, price_b float64, pair_rate map[string]float64,
	order_buy getinfobinance.AdvertiserAdv, paramUser interact.Parameters) result.ResultP2P {

	profitResult := result.ResultP2P{}

	var transAmountSecond float64
	//log.Printf("Pair rate %s - %v", p, pair_rate[p])
	var assetSell string
	if strings.HasPrefix(p, a) {
		transAmountSecond = (transAmountFirst * pair_rate[p])
		assetSell = p[len(a):]
		//log.Println(convertfiat, transAmountSecond)
	} else {
		transAmountSecond = (transAmountFirst / pair_rate[p])
		assetSell = p[:(len(p) - len(a))]
		//log.Println(convertfiat, transAmountSecond)
	}
	//third steps
	order_sell := getinfobinance.GetDataP2P(assetSell, fiat,
		"Sell", paramUser)
	//log.Printf("First - %s, Second - %s", a, convertfiat)
	if order_sell.Adv.Price == 0 {
		return profitResult
	}
	price_s := order_sell.Adv.Price

	transAmountThird := price_s * transAmountSecond
	//log.Printf("transAmountFirst - %v, transAmountSecond - %v, transAmountThird - %v",
	//	transAmountFirst, transAmountSecond, transAmountThird)
	//
	transAmountFloat, err := strconv.ParseFloat(paramUser.TransAmount, 64)
	if err != nil {
		log.Printf("Problem with convert transAmount to float, err - %v", err)
	}
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
	return profitResult
}
