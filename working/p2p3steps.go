package working

import (
	"encoding/json"
	"fmt"
	"github.com/Zmey56/arbitrage/getdata"
	"github.com/Zmey56/arbitrage/getinfobinance"
	"github.com/Zmey56/arbitrage/interact"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

type ResultP2P struct {
	Profit          bool
	DataTime        time.Time
	Fiat            string
	AssetsBuy       string
	PriceAssetsBuy  float64
	LinkAssetsBuy   string
	Pair            string
	PricePair       float64
	LinkMarket      string
	AssetsSell      string
	PriceAssetsSell float64
	LinkAssetsSell  string
	ProfitValue     float64
	ProfitPercet    string
}

func P2P3steps(fiat string, paramUser interact.Parameters) {
	allOrders := [][]ResultP2P{}
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
			arr_val := GetResultP2P3(a, fiat, pair, paramUser)
			allOrders = append(allOrders, arr_val)
		}(a)
	}
	wg.Wait()

	for _, j := range allOrders {
		for _, i := range j {
			if i.Profit {
				saveJsonFile(fiat, i, true)
			} else {
				saveJsonFile(fiat, i, false)
			}
		}
	}
}

func GetResultP2P3(a, fiat string, pair map[string][]string, paramUser interact.Parameters) []ResultP2P {
	//fmt.Println("====================================")
	//log.Println("ASSETS", a)
	var resultP2PArr []ResultP2P
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
				value := PrintResultP2P3(p, a, fiat, transAmountFirst, price_b,
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

func PrintResultP2P3(p, a, fiat string, transAmountFirst, price_b float64, pair_rate map[string]float64,
	order_buy getinfobinance.AdvertiserAdv, paramUser interact.Parameters) ResultP2P {

	profitResult := ResultP2P{}

	var transAmountSecond float64
	//log.Printf("Pair rate %s - %v", p, pair_rate[p])
	var assetSell string
	if strings.HasPrefix(p, a) {
		transAmountSecond = (transAmountFirst * pair_rate[p])
		assetSell = p[len(a):]
		//log.Println(convertfiat, transAmountSecond)
	} else {
		transAmountSecond = (transAmountFirst / pair_rate[p])
		assetSell = p[:(len(p) - len(a))]
		//log.Println(convertfiat, transAmountSecond)
	}
	//third steps
	order_sell := getinfobinance.GetDataP2P(assetSell, fiat,
		"Sell", paramUser)
	//log.Printf("First - %s, Second - %s", a, convertfiat)
	if order_sell.Adv.Price == "" {
		return profitResult
	}
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
	profitResult.Profit = transAmountThird > transAmountFloat
	profitResult.DataTime = time.Now()
	profitResult.Fiat = fiat
	profitResult.AssetsBuy = a
	profitResult.PriceAssetsBuy = price_b
	profitResult.LinkAssetsBuy = fmt.Sprintf("https://p2p.binance.com/en/advertiserDetail?advertiserNo=%v", order_buy.Adv.AdvNo)
	profitResult.Pair = p
	profitResult.PricePair = pair_rate[p]
	profitResult.LinkMarket = returnLinkMarket(a, p)
	profitResult.AssetsSell = assetSell
	profitResult.PriceAssetsSell = price_s
	profitResult.LinkAssetsSell = fmt.Sprintf("https://p2p.binance.com/en/advertiserDetail?advertiserNo=%v", order_sell.Adv.AdvNo)
	profitResult.ProfitValue = transAmountThird - transAmountFloat
	profitResult.ProfitPercet = fmt.Sprintf("%.2f", ((transAmountThird-transAmountFloat)/transAmountFloat)*100)
	return profitResult

}

func returnLinkMarket(a, p string) string {
	var pair string
	if strings.HasPrefix(p, a) {
		pair = a + "_" + p[len(a):]
	} else {
		pair = p[:(len(p)-len(a))] + "_" + a
	}
	return fmt.Sprintf("https://www.binance.com/en/trade/%v?_from=markets", pair)
}

func saveJsonFile(fiat string, pr ResultP2P, profit bool) {
	current := time.Now()
	current.String()
	var path_save string
	if profit {
		path_save = fmt.Sprintf("jsonresult/POSITIVE%s_%s.json", fiat, current.Format("2006_01_02"))
	} else {
		path_save = fmt.Sprintf("jsonresult/NEGATIVE%s_%s.json", fiat, current.Format("2006_01_02"))
	}
	tmp_result := []ResultP2P{}
	if exists(path_save) {
		jsonfile, err := os.ReadFile(path_save)
		if err != nil {
			panic(err)
		}

		_ = json.Unmarshal(jsonfile, &tmp_result)
	}
	tmp_result = append(tmp_result, pr)
	f, err := os.OpenFile(path_save,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
	}
	file, _ := json.MarshalIndent(tmp_result, "", " ")
	os.WriteFile(path_save, file, 0666)
	defer f.Close()
}

func exists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}
