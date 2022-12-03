package working

import (
	"fmt"
	"github.com/Zmey56/arbitrage/getdata"
	"github.com/Zmey56/arbitrage/getinfobinance"
	"log"
	"strings"
	"sync"
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
	var wg sync.WaitGroup
	for _, a := range assets_symbol {
		wg.Add(1)
		go func(a string) {
			defer wg.Done()
			GetResultP2P3(a, fiat, transAmount, pair, payTypes)
		}(a)
	}
	wg.Wait()
}

func GetResultP2P3(a, fiat string, transAmount float64, pair map[string][]string, payTypes []string) {
	//fmt.Println("====================================")
	//log.Println("ASSETS", a)
	pair_assets := pair[a]
	//first step
	order_buy, price_b := getinfobinance.GetDataP2PVer2(a, fiat, "Buy", transAmount, payTypes)
	transAmountFirst := transAmount / price_b
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
			PrintResultP2P3(p, a, fiat, transAmount, transAmountFirst, price_b,
				pair_rate, order_buy, payTypes)
		}(p)

	}
	wg.Wait()
}

func PrintResultP2P3(p, a, fiat string, transAmount, transAmountFirst, price_b float64,
	pair_rate map[string]float64, order_buy getinfobinance.AdvertiserAdv, payTypes []string) {
	var transAmountSecond float64
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
	order_sell, price_s := getinfobinance.GetDataP2PVer2(convertfiat, fiat,
		"Sell", transAmountSecond, payTypes)
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
			fmt.Sprintf("%.2f", ((transAmountThird-transAmount)/transAmount)*100), "\n")
		fmt.Println("Payment Methods Buy: ", order_buy.Adv.TradeMethods[0].TradeMethodShortName)

		fmt.Println("Payment Methods Buy: ", order_sell.Adv.TradeMethods[0].TradeMethodShortName)
		fmt.Println("\n")
	}
}
