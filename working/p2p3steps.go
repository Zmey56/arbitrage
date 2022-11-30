package working

import (
	"fmt"
	"github.com/Zmey56/arbitrage/getdata"
	"github.com/Zmey56/arbitrage/getinfobinance"
	"log"
	"strconv"
	"strings"
)

func P2P3steps(fiat string, payTypes []string, transAmount float64) {
	//get all assets from binance for this fiat
	assets := getinfobinance.GetAssets(fiat)
	assets_symbol := make([]string, 0, len(assets))
	assets_name := make([]string, 0, len(assets))

	for k, v := range assets {
		assets_symbol = append(assets_symbol, k)
		assets_name = append(assets_name, v)
	}

	//get pair for rate

	pair := getdata.GetPairFromJSON(fiat)

	//get information about orders with binance
	for _, a := range assets_symbol {
		//fmt.Println("====================================")
		//log.Println("ASSETS", a)
		pair_assets := pair[a]
		//first step
		order_buy := getinfobinance.GetDataP2P(a, fiat, "Buy", transAmount, payTypes)
		price_b, _ := strconv.ParseFloat(order_buy.Advs.Price, 64)
		transAmountFirst := transAmount / price_b
		//log.Printf("transAmountFirst - %v", transAmountFirst)
		//second step
		var transAmountSecond float64
		//log.Printf("Pair Assets - %s", pair_assets)
		pair_rate := getinfobinance.GetRatePair(pair_assets)
		//log.Printf("Pair Rate - %s", pair_rate)
		for p := range pair_rate {
			//log.Printf("Pair rate %s - %v", p, pair_rate[p])
			var convertfiat string
			if strings.HasPrefix(p, a) {
				transAmountSecond = (transAmountFirst * pair_rate[p])
				convertfiat = p[len(a):]
				//log.Println(convertfiat, transAmountSecond)
			} else {
				transAmountSecond = (transAmountFirst / pair_rate[p])
				convertfiat = p[:(len(p) - len(a))]
				//log.Println(convertfiat, transAmountSecond)
			}
			//third steps
			order_sell := getinfobinance.GetDataP2P(convertfiat, fiat, "Sell", transAmountSecond, payTypes)
			price_s, _ := strconv.ParseFloat(order_sell.Advs.Price, 64)
			//log.Printf("First - %s, Second - %s", a, convertfiat)
			transAmountThird := price_s * transAmountSecond
			//log.Printf("transAmountFirst - %v, transAmountSecond - %v, transAmountThird - %v",
			//	transAmountFirst, transAmountSecond, transAmountThird)
			//
			if transAmountThird > transAmount {
				log.Printf("assets - %s, pair - %s, convertfiat - %s, fiat - %s",
					a, p, convertfiat, fiat)
				log.Printf("PriceFirst - %v, PriceSecond - %v, Price Third - %v",
					price_b, pair_rate[p], price_s)

				fmt.Println("RESULT:", fmt.Sprintf("%.2f", transAmountThird-transAmount), " ",
					fmt.Sprintf("%.2f", ((transAmountThird-transAmount)/transAmount)*100))
			}
		}
	}
}
