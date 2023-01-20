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

func P2P3stepsTakerMaker(fiat string, paramUser workingbinance.ParametersBinance) {
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
			arr_val := GetResultP2P3TM(a, fiat, pair, paramUser)
			allOrders = append(allOrders, arr_val)
		}(a)
	}
	wg.Wait()

	for _, j := range allOrders {
		for _, i := range j {
			if (i.TotalAdvBuy > 0) && (i.TotalAdvSell > 0) {
				result.SaveResultJsonFile(fiat, i, "3steps_tm")
				//log.Printf("3 steps taker maker. Fiat - %s, Result - %v", fiat, i)
				if i.Profit && (i.ProfitPercet >= paramUser.PercentUser) && (i.TotalAdvBuy > 10) && (i.TotalAdvSell > 10) {
					result.FormatMessageAndSend(i, "You are Taker", "You are Maker")
				}
			}
		}
	}
}

func GetResultP2P3TM(a, fiat string, pair map[string][]string, paramUser workingbinance.ParametersBinance) []result.ResultP2P {

	var resultP2PArr []result.ResultP2P
	pair_assets := pair[a]
	//first step
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
		log.Println("New transAmount because didn't enter amount in beginer", paramUser.TransAmount)
	}

	if len(order_buy.Data) > 0 {

		price_b := order_buy.Data[0].Adv.Price

		transAmountFirst := transAmountFloat / price_b
		//second step

		pair_rate := getinfobinance.GetRatePair(pair_assets)

		var wg sync.WaitGroup
		for p := range pair_rate {
			wg.Add(1)

			go func(p string) {
				defer wg.Done()
				value := PrintResultP2P3TM(p, a, fiat, transAmountFirst, price_b,
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

func PrintResultP2P3TM(p, a, fiat string, transAmountFirst, price_b float64, pair_rate map[string]float64,
	order_buy getinfobinance.AdvertiserAdv, paramUser workingbinance.ParametersBinance) result.ResultP2P {

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
		"Buy", paramUser)
	if len(order_sell.Data) == 0 {
		return profitResult
	}
	price_s := order_sell.Data[0].Adv.Price

	transAmountThird := price_s * transAmountSecond

	transAmountFloat, err := strconv.ParseFloat(paramUser.TransAmount, 64)
	if err != nil {
		log.Printf("Problem with convert transAmount to float, err - %v", err)
	}
	profitResult.Market.First = "Binance"
	profitResult.Market.Second = "Binance"
	profitResult.Market.Third = "Binance"
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
	profitResult.LinkAssetsSell = fmt.Sprintf("https://p2p.binance.com/en/trade/all-payments/%v?fiat=%v",
		assetSell, fiat)
	profitResult.ProfitValue = transAmountThird - transAmountFloat
	profitResult.ProfitPercet = (((transAmountThird - transAmountFloat) / transAmountFloat) * 100)
	profitResult.TotalAdvBuy = order_buy.Total
	profitResult.TotalAdvSell = order_sell.Total
	profitResult.AdvNoBuy = order_buy.Data[0].Adv.AdvNo
	profitResult.AdvNoSell = order_sell.Data[0].Adv.AdvNo
	return profitResult
}
