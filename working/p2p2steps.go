// function for get info about p2p 2 steps
package working

import (
	"fmt"
	"github.com/Zmey56/arbitrage/getinfobinance"
	"strconv"
)

func P2P2steps(fiat string, payTypes []string, transAmount float64) {
	//get all assets from binance for this fiat
	assets := getinfobinance.GetAssets(fiat)
	assets_symbol := make([]string, 0, len(assets))
	assets_name := make([]string, 0, len(assets))

	for k, v := range assets {
		assets_symbol = append(assets_symbol, k)
		assets_name = append(assets_name, v)
	}

	//pair := getdata.GetPairFromJSON(fiat)

	//get information about orders with binance
	for _, a := range assets_symbol {
		fmt.Println("====================================")
		fmt.Println("ASSETS", a)
		order_buy := getinfobinance.GetDataP2P(a, fiat, "Buy", transAmount, payTypes)
		order_sell := getinfobinance.GetDataP2P(a, fiat, "Sell", transAmount, payTypes)
		price_b, _ := strconv.ParseFloat(order_buy.Advs.Price, 64)
		fmt.Println("BUY ", price_b)
		fmt.Println("Nick", order_buy.Advertisers.NickName)
		fmt.Println("Orders", order_buy.Advertisers.MonthOrderCount)
		rate := fmt.Sprintf("%.2f", order_buy.Advertisers.MonthFinishRate*100.00)
		fmt.Println("Completition", rate, "%")
		fmt.Println("Available", order_buy.Advs.TradableQuantity)
		fmt.Println("Limit", order_buy.Advs.MinSingleTransAmount, " - ", order_buy.Advs.MaxSingleTransAmount)
		fmt.Println("Payments", order_buy.Advs.TradeMethods.TradeMethodName)
		fmt.Println("------------------------------------")
		price_s, _ := strconv.ParseFloat(order_sell.Advs.Price, 64)
		fmt.Println("SELL ", price_s)
		fmt.Println("Nick", order_sell.Advertisers.NickName)
		fmt.Println("Orders", order_sell.Advertisers.MonthOrderCount)
		rate = fmt.Sprintf("%.2f", order_sell.Advertisers.MonthFinishRate*100.00)
		fmt.Println("Completition", rate, "%")
		fmt.Println("Available", order_sell.Advs.TradableQuantity)
		fmt.Println("Limit", order_sell.Advs.MinSingleTransAmount, " - ", order_sell.Advs.MaxSingleTransAmount)
		fmt.Println("Payments", order_sell.Advs.TradeMethods.TradeMethodName)
		fmt.Println("\n")
		fmt.Println("RESULT:", fmt.Sprintf("%.2f", price_s-price_b), " ",
			fmt.Sprintf("%.2f", ((price_s-price_b)/price_b)*100))
		fmt.Println("====================================")
	}

}
