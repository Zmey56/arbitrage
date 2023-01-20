// function for get info about p2p 2 steps
package p2pbinance

import (
	"fmt"
	"github.com/Zmey56/arbitrage/pkg/getdata"
	"github.com/Zmey56/arbitrage/pkg/workingbinance"
)

func P2P2stepsBinance(fiat string, paramUser workingbinance.ParametersBinance) {
	//get all assets from binance for this fiat
	assets := getdata.GetAssets(fiat)
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
		order_buy := getdata.GetDataP2PBinance(a, fiat, "Buy", paramUser)
		order_sell := getdata.GetDataP2PBinance(a, fiat, "Sell", paramUser)
		price_b := order_buy.Data[0].Adv.Price
		price_s := order_sell.Data[0].Adv.Price
		fmt.Println("BUY ", price_b)
		fmt.Println("Nick", order_buy.Data[0].Advertiser.NickName)
		fmt.Println("Orders", order_buy.Data[0].Advertiser.MonthOrderCount)
		rate := fmt.Sprintf("%.2f", order_buy.Data[0].Advertiser.MonthFinishRate*100.00)
		fmt.Println("Completition", rate, "%")
		fmt.Println("Available", order_buy.Data[0].Adv.TradableQuantity)
		fmt.Println("Limit", order_buy.Data[0].Adv.MinSingleTransAmount, " - ", order_buy.Data[0].Adv.MaxSingleTransAmount)
		fmt.Println("Payments", order_buy.Data[0].Adv.TradeMethods)
		fmt.Println("------------------------------------")
		fmt.Println("SELL ", price_s)
		fmt.Println("Nick", order_sell.Data[0].Advertiser.NickName)
		fmt.Println("Orders", order_sell.Data[0].Advertiser.MonthOrderCount)
		rate = fmt.Sprintf("%.2f", order_sell.Data[0].Advertiser.MonthFinishRate*100.00)
		fmt.Println("Completition", rate, "%")
		fmt.Println("Available", order_sell.Data[0].Adv.TradableQuantity)
		fmt.Println("Limit", order_sell.Data[0].Adv.MinSingleTransAmount, " - ", order_sell.Data[0].Adv.MaxSingleTransAmount)
		fmt.Println("Payments", order_sell.Data[0].Adv.TradeMethods)
		fmt.Println("\n")
		fmt.Println("RESULT:", fmt.Sprintf("%.2f", price_s-price_b), " ",
			fmt.Sprintf("%.2f", ((price_s-price_b)/price_b)*100))
		fmt.Println("====================================")
	}
}
