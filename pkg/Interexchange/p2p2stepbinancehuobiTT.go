package Interexchange

import (
	"fmt"
	"github.com/Zmey56/arbitrage/pkg/getdata"
	"github.com/Zmey56/arbitrage/pkg/getdatahuobi"
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

func P2P2stepsBinanceHuobiTT(fiat string, binance workingbinance.ParametersBinance, huobi getinfohuobi.ParametersHuobi) {
	allOrders := []result.ResultP2P{}
	//get all assets from binance for this fiat

	boarder := 5
	if fiat == "RUS" {
		boarder = 10
	}
	assets := getdata.GetAssets(fiat)
	assets_symbol := make([]string, 0, len(assets))
	assets_name := make([]string, 0, len(assets))

	for k, v := range assets {
		assets_symbol = append(assets_symbol, k)
		assets_name = append(assets_name, v)
	}

	//get information about orders with binance
	var wg sync.WaitGroup
	for _, a := range assets_symbol {
		wg.Add(1)
		go func(a string) {
			defer wg.Done()
			arr_val := getResultP2P2TTBH(a, fiat, binance, huobi)
			allOrders = append(allOrders, arr_val)
		}(a)

		wg.Add(1)
		go func(a string) {
			defer wg.Done()
			arr_val := getResultP2P2TTHB(a, fiat, binance, huobi)
			allOrders = append(allOrders, arr_val)
		}(a)
	}
	wg.Wait()

	for _, j := range allOrders {
		if j.TotalAdvBuy > 0 && j.TotalAdvSell > 0 {
			result.SaveResultJsonFile(fiat, j, "2stepsTT")
			//log.Printf("2 steps taker taker. Fiat - %s, Result - %v", fiat, j)
			if j.Profit && (j.ProfitPercet >= binance.PercentUser) && (j.TotalAdvBuy >= boarder) && (j.TotalAdvSell >= boarder) {
				result.FormatMessageAndSend2steps(j, "You are Taker", "You are Taker")
			}
		}
	}
}

func getResultP2P2TTBH(a, fiat string, binance workingbinance.ParametersBinance,
	huobi getinfohuobi.ParametersHuobi) result.ResultP2P {

	coinidmap := workinghuobi.GetCoinIDHuobo(fiat)

	//log.Println(coinidmap[a])

	profitResult := result.ResultP2P{}

	order_buy := getdata.GetDataP2PBinance(a, fiat, "Buy", binance)

	if len(order_buy.Data) > 0 {
		var transAmountThirst, volume float64

		if binance.TransAmount != "" {
			huobi.Amount = binance.TransAmount
			transAmountThirst, _ = strconv.ParseFloat(binance.TransAmount, 64)
			volume = transAmountThirst / order_buy.Data[0].Adv.Price
		} else {
			transAmountThirst, _ = strconv.ParseFloat(order_buy.Data[0].Adv.DynamicMaxSingleTransAmount, 64)
			huobi.Amount = order_buy.Data[0].Adv.DynamicMaxSingleTransAmount
		}

		order_sell := getdatahuobi.GetDataP2PHuobi(coinidmap[strings.ToUpper(a)], coinidmap[fiat],
			"buy", huobi)
		if len(order_sell.Data) == 0 {
			return profitResult
		}
		price_b := order_buy.Data[0].Adv.Price
		price_s, _ := strconv.ParseFloat(order_sell.Data[0].Price, 64)

		transAmountSecond := price_s * volume
		log.Println("FIAT", fiat, "COIN", a, "price_b", price_b, "price_s", price_s, "\n")

		profitResult.Market.First = "binance"
		profitResult.Market.Second = "huobi"
		profitResult.Profit = transAmountSecond > transAmountThirst
		profitResult.DataTime = time.Now()
		profitResult.Fiat = fiat
		profitResult.AssetsBuy = a
		profitResult.PriceAssetsBuy = price_b
		profitResult.PaymentBuy = result.PaymentMetods(order_buy)
		profitResult.LinkAssetsBuy = fmt.Sprintf("https://p2p.binance.com/en/trade/all-payments/%v?fiat=%v", a, fiat)
		profitResult.AssetsSell = strings.ToUpper(a)
		profitResult.PriceAssetsSell = price_s
		profitResult.PaymentSell = result.PaymentMetodsHuobi(order_sell)
		profitResult.LinkAssetsSell = fmt.Sprintf("https://www.huobi.com/en-us/fiat-crypto/trader/%s", strconv.Itoa(order_sell.Data[0].UID))
		profitResult.ProfitValue = transAmountSecond - transAmountThirst
		profitResult.ProfitPercet = (((transAmountSecond - transAmountThirst) / transAmountThirst) * 100)
		profitResult.TotalAdvBuy = order_buy.Total
		profitResult.TotalAdvSell = order_sell.TotalCount
		profitResult.AdvNoBuy = order_buy.Data[0].Adv.AdvNo
		profitResult.AdvNoSell = strconv.Itoa(order_sell.Data[0].UID)
		//log.Println("RESULT:", profitResult)
		//fmt.Printf("%s - %s %+v\n", fiat, a, profitResult)
		return profitResult
	} else {
		return profitResult
	}

}

func getResultP2P2TTHB(a, fiat string, binance workingbinance.ParametersBinance,
	huobi getinfohuobi.ParametersHuobi) result.ResultP2P {

	coinidmap := workinghuobi.GetCoinIDHuobo(fiat)

	//log.Println(coinidmap[a])

	profitResult := result.ResultP2P{}

	order_buy := getdatahuobi.GetDataP2PHuobi(coinidmap[a], coinidmap[fiat], "sell", huobi)

	if len(order_buy.Data) > 0 {
		var transAmountThirst, volume float64
		price_b, _ := strconv.ParseFloat(order_buy.Data[0].Price, 64)

		if huobi.Amount != "" {
			binance.TransAmount = huobi.Amount
			transAmountThirst, _ = strconv.ParseFloat(huobi.Amount, 64)
			volume = transAmountThirst / price_b
		} else {
			transAmountThirst, _ = strconv.ParseFloat(order_buy.Data[0].MaxTradeLimit, 64)
			binance.TransAmount = order_buy.Data[0].MaxTradeLimit
		}

		order_sell := getdata.GetDataP2PBinance(a, fiat, "sell", binance)
		if len(order_sell.Data) == 0 {
			return profitResult
		}

		price_s := order_sell.Data[0].Adv.Price

		transAmountSecond := price_s * volume

		profitResult.Market.First = "huobi"
		profitResult.Market.Second = "binance"
		profitResult.Profit = transAmountSecond > transAmountThirst
		profitResult.DataTime = time.Now()
		profitResult.Fiat = fiat
		profitResult.AssetsBuy = a
		profitResult.PriceAssetsBuy = price_b
		profitResult.PaymentBuy = result.PaymentMetodsHuobi(order_buy)
		profitResult.LinkAssetsBuy = fmt.Sprintf("https://www.huobi.com/en-us/fiat-crypto/trade/buy-%s-%s", a, fiat)
		profitResult.AssetsSell = strings.ToUpper(a)
		profitResult.PriceAssetsSell = price_s
		profitResult.PaymentSell = result.PaymentMetods(order_sell)
		profitResult.LinkAssetsSell = fmt.Sprintf("https://p2p.binance.com/en/trade/sell/%v?fiat=%v", a, fiat)
		profitResult.ProfitValue = transAmountSecond - transAmountThirst
		profitResult.ProfitPercet = (((transAmountSecond - transAmountThirst) / transAmountThirst) * 100)
		profitResult.TotalAdvSell = order_buy.TotalCount
		profitResult.TotalAdvBuy = order_sell.Total
		profitResult.AdvNoBuy = strconv.Itoa(order_buy.Data[0].UID)
		profitResult.AdvNoSell = order_sell.Data[0].Adv.AdvNo
		//log.Println("RESULT:", profitResult)
		//fmt.Printf("%s - %s %+v\n", fiat, a, profitResult)
		return profitResult
	} else {
		return profitResult
	}
}
