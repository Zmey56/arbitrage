// function for get info about p2p 2 steps
package working

import (
	"fmt"
	"github.com/Zmey56/arbitrage/getdata"
	"github.com/Zmey56/arbitrage/getinfobinance"
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

	//get information about orders with binance
	for _, a := range assets_symbol {
		order_buy := getdata.GetDataP2P(a, fiat, "Buy", transAmount, payTypes)
		order_sell := getdata.GetDataP2P(a, fiat, "Sell", transAmount, payTypes)
		fmt.Println("BUY")
		fmt.Println(order_buy)
		fmt.Println("SELL")
		fmt.Println(order_sell)
	}

}
