package p2phuobi

import (
	"fmt"
	"github.com/Zmey56/arbitrage/pkg/getdatahuobi"
	"github.com/Zmey56/arbitrage/pkg/getinfohuobi"
	"github.com/Zmey56/arbitrage/pkg/workinghuobi"
	"log"
	"strconv"
)

func P2P2stepsHuobi(fiat string, paramUser getinfohuobi.ParametersHuobi) {
	coinidmap := workinghuobi.GetCoinIDHuobo(fiat)
	log.Println(coinidmap)

	//get all assets from binance for this fiat
	currencyarr := getdatahuobi.GetCurrencyHuobi(fiat)
	log.Println(currencyarr)

	//pair := getdata.GetPairFromJSON(fiat)

	//get information about orders with binance
	for _, a := range currencyarr {
		fmt.Println("====================================")
		fmt.Println("ASSETS", a)
		log.Println(fiat, " - ", coinidmap[fiat], " | ", a, " - ", coinidmap[a])
		order_buy := getdatahuobi.GetDataP2PHuobi(coinidmap[a], coinidmap[fiat], "sell", paramUser)
		order_sell := getdatahuobi.GetDataP2PHuobi(coinidmap[a], coinidmap[fiat], "buy", paramUser)
		if len(order_buy.Data) > 0 || len(order_sell.Data) > 0 {
			price_b, _ := strconv.ParseFloat(order_buy.Data[0].Price, 64)
			price_s, _ := strconv.ParseFloat(order_sell.Data[0].Price, 64)
			fmt.Println("BUY ", price_b)
			fmt.Println("Nick", order_buy.Data[0].UserName)
			fmt.Println("Orders", order_buy.Data[0].TradeMonthTimes)
			rate := fmt.Sprintf(order_buy.Data[0].OrderCompleteRate)
			fmt.Println("Completition", rate, "%")

			fmt.Println("Available", order_buy.Data[0].TradeCount)
			fmt.Println("Limit", order_buy.Data[0].MinTradeLimit, " - ", order_buy.Data[0].MaxTradeLimit)
			fmt.Println("Payments", order_buy.Data[0].PayMethods)
			fmt.Println("------------------------------------")
			fmt.Println("SELL ", price_s)
			fmt.Println("Nick", order_sell.Data[0].UserName)
			fmt.Println("Orders", order_sell.Data[0].TradeCount)
			rate = fmt.Sprintf("%.2f", order_sell.Data[0].TradeMonthTimes)
			fmt.Println("Completition", rate, "%")
			fmt.Println("Available", order_sell.Data[0].TradeCount)
			fmt.Println("Limit", order_sell.Data[0].MinTradeLimit, " - ", order_sell.Data[0].MaxTradeLimit)
			fmt.Println("Payments", order_sell.Data[0].PayMethods)
			fmt.Println("\n")
			fmt.Println("RESULT:", fmt.Sprintf("%.2f", price_s-price_b), " ",
				fmt.Sprintf("%.2f", ((price_s-price_b)/price_b)*100))
			fmt.Println("====================================")
		}
	}

}
