package working

import (
	"fmt"
	"github.com/Zmey56/arbitrage/Interact"
	"github.com/Zmey56/arbitrage/getdata"
	"github.com/Zmey56/arbitrage/getinfobinance"
	"log"
	"strconv"
	"strings"
	"sync"
)

func P2P3steps(fiat string, paramUser Interact.Parameters) {
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
			GetResultP2P3(a, fiat, pair, paramUser)
		}(a)
	}
	wg.Wait()
}

func GetResultP2P3(a, fiat string, pair map[string][]string, paramUser Interact.Parameters) {
	//fmt.Println("====================================")
	//log.Println("ASSETS", a)
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
	if order_buy.Adv.Price != "" {
		price_b, err := strconv.ParseFloat(order_buy.Adv.Price, 64)
		if err != nil {
			log.Printf("Can't parse string to float for price buy, error: %s", err)
		}
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
				PrintResultP2P3(p, a, fiat, transAmountFirst, price_b,
					pair_rate, order_buy, paramUser)
			}(p)

		}
		wg.Wait()
	} else {
		return
	}

}

func PrintResultP2P3(p, a, fiat string, transAmountFirst, price_b float64,
	pair_rate map[string]float64, order_buy getinfobinance.AdvertiserAdv, paramUser Interact.Parameters) {

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
	order_sell := getinfobinance.GetDataP2P(convertfiat, fiat,
		"Sell", paramUser)
	//log.Printf("First - %s, Second - %s", a, convertfiat)
	price_s, err := strconv.ParseFloat(order_sell.Adv.Price, 64)
	if err != nil {
		log.Printf("Can't parse string to float for price sell, error: %s", err)
	}
	transAmountThird := price_s * transAmountSecond
	//log.Printf("transAmountFirst - %v, transAmountSecond - %v, transAmountThird - %v",
	//	transAmountFirst, transAmountSecond, transAmountThird)
	//
	transAmountFloat, err := strconv.ParseFloat(paramUser.TransAmount, 64)
	if err != nil {
		log.Printf("Problem with convert transAmount to float, err - %v", err)
	}

	//log.Printf("transAmoountStart - %v and transAmoountEnd - %v", paramUser.TransAmount, transAmountThird)
	//log.Printf("assets - %s, pair - %s, convertfiat - %s, fiat - %s",
	//	a, p, convertfiat, fiat)
	//log.Printf("PriceFirst - %v, PriceSecond - %v, Price Third - %v",
	//	price_b, pair_rate[p], price_s)
	//log.Println("RESULT:", fmt.Sprintf("%.2f", transAmountThird-transAmountFloat), " ",
	//	fmt.Sprintf("%.2f", ((transAmountThird-transAmountFloat)/transAmountFloat)*100), "\n")
	//log.Println("Payment Methods Buy: ", order_buy.Adv.TradeMethods[0].TradeMethodShortName)
	//log.Printf("User Name - %s, orders - %v, completion - %v",
	//	order_buy.Advertiser.NickName, order_buy.Advertiser.MonthOrderCount,
	//	order_buy.Advertiser.MonthFinishRate*100)
	//log.Printf("Link: https://p2p.binance.com/en/advertiserDetail?advertiserNo=%s", order_buy.Advertiser.UserNo)
	//log.Println("Payment Methods Sell: ", order_sell.Adv.TradeMethods[0].TradeMethodShortName)
	//log.Printf("User Name - %s, orders - %v, completion - %v",
	//	order_sell.Advertiser.NickName, order_sell.Advertiser.MonthOrderCount,
	//	order_sell.Advertiser.MonthFinishRate*100)
	//log.Printf("Link: https://p2p.binance.com/en/advertiserDetail?advertiserNo=%s", order_sell.Advertiser.UserNo)
	//log.Println("\n")

	if transAmountThird > transAmountFloat {
		fmt.Println("!!!!!!!!!NB!!!!!!!!!!!!!")
		fmt.Printf("assets - %s, pair - %s, convertfiat - %s, fiat - %s\n",
			a, p, convertfiat, fiat)
		fmt.Printf("PriceFirst - %v, PriceSecond - %v, Price Third - %v\n",
			price_b, pair_rate[p], price_s)
		fmt.Println("RESULT:", fmt.Sprintf("%.2f", transAmountThird-transAmountFloat), " ",
			fmt.Sprintf("%.2f", ((transAmountThird-transAmountFloat)/transAmountFloat)*100), "\n")
		fmt.Println("Payment Methods Buy: ", order_buy.Adv.TradeMethods[0].TradeMethodShortName)

		fmt.Println("Payment Methods Sell: ", order_sell.Adv.TradeMethods[0].TradeMethodShortName)
		fmt.Println("\n")
	}
}
