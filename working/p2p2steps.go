// function for get info about p2p 2 steps
package working

import (
	"fmt"
	"github.com/Zmey56/arbitrage/Interact"
	"github.com/Zmey56/arbitrage/getinfobinance"
	"log"
	"strconv"
)

func P2P2steps(fiat string, paramUser Interact.Parameters) {
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
		order_buy := getinfobinance.GetDataP2P(a, fiat, "Buy", paramUser)
		order_sell := getinfobinance.GetDataP2P(a, fiat, "Sell", paramUser)
		price_b, err := strconv.ParseFloat(order_buy.Adv.Price, 64)
		if err != nil {
			log.Printf("Can't parse string to float for price buy, error: %s", err)
		}
		price_s, err := strconv.ParseFloat(order_sell.Adv.Price, 64)
		if err != nil {
			log.Printf("Can't parse string to float for price sell, error: %s", err)
		}
		fmt.Println("BUY ", price_b)
		fmt.Println("Nick", order_buy.Advertiser.NickName)
		fmt.Println("Orders", order_buy.Advertiser.MonthOrderCount)
		rate := fmt.Sprintf("%.2f", order_buy.Advertiser.MonthFinishRate*100.00)
		fmt.Println("Completition", rate, "%")
		fmt.Println("Available", order_buy.Adv.TradableQuantity)
		fmt.Println("Limit", order_buy.Adv.MinSingleTransAmount, " - ", order_buy.Adv.MaxSingleTransAmount)
		fmt.Println("Payments", order_buy.Adv.TradeMethods)
		fmt.Println("------------------------------------")
		fmt.Println("SELL ", price_s)
		fmt.Println("Nick", order_sell.Advertiser.NickName)
		fmt.Println("Orders", order_sell.Advertiser.MonthOrderCount)
		rate = fmt.Sprintf("%.2f", order_sell.Advertiser.MonthFinishRate*100.00)
		fmt.Println("Completition", rate, "%")
		fmt.Println("Available", order_sell.Adv.TradableQuantity)
		fmt.Println("Limit", order_sell.Adv.MinSingleTransAmount, " - ", order_sell.Adv.MaxSingleTransAmount)
		fmt.Println("Payments", order_sell.Adv.TradeMethods)
		fmt.Println("\n")
		fmt.Println("RESULT:", fmt.Sprintf("%.2f", price_s-price_b), " ",
			fmt.Sprintf("%.2f", ((price_s-price_b)/price_b)*100))
		fmt.Println("====================================")
	}

}
