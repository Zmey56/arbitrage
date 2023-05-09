package withoutcard

import (
	"encoding/json"
	"fmt"
	"github.com/Zmey56/arbitrage/pkg/commonfunction"
	"github.com/Zmey56/arbitrage/pkg/getdata"
	"github.com/Zmey56/arbitrage/pkg/getinfohuobi"
	"github.com/Zmey56/arbitrage/pkg/getinfookx"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

var urlBinanceExchangeInfo = "https://api.binance.com/api/v3/exchangeInfo" //about all pair Binance
var urlHuobiExchangeInfo = "https://api.huobi.pro/v1/common/symbols"       //about all pair Huobi

// GetPairHuobi - get all value pair-pair for use binance-pairHuobi-pairHuobi

func GetPairBHH(asset string) {
	//get all pair for asset from Binance
	respB, err := http.Get(urlBinanceExchangeInfo)
	if err != nil {
		panic(err)
	}
	defer respB.Body.Close()

	bodyB, _ := io.ReadAll(respB.Body)

	symbolsB := getdata.Symbols{}

	if err := json.Unmarshal(bodyB, &symbolsB); err != nil {
		panic(err)
	}

	allPairAssetB := []string{}

	for _, valueB := range symbolsB.Symbols {
		if valueB.BaseAsset == asset || valueB.QuoteAsset == asset {
			allPairAssetB = append(allPairAssetB, valueB.Symbol)
		}

	}

	//get all pair for asset from Huobi for end

	respH, err := http.Get(urlHuobiExchangeInfo)
	if err != nil {
		panic(err)
	}
	defer respH.Body.Close()

	bodyH, _ := io.ReadAll(respH.Body)

	symbolsH := getinfohuobi.PairHuobi{}

	if err := json.Unmarshal(bodyH, &symbolsH); err != nil {
		panic(err)
	}

	allPairAssetH := []string{}
	allPairH := []string{} // get all real pair from Huobi

	for _, valueH := range symbolsH.Data {
		allPairH = append(allPairH, valueH.Symbol)
		if valueH.BaseCurrency == strings.ToLower(asset) || valueH.QuoteCurrency == strings.ToLower(asset) {
			allPairAssetH = append(allPairAssetH, valueH.Symbol)
		}
	}

	// forming map for binance-pair-pair
	mapPairFinal := make(map[string][]string)

	// remove asset from first pair
	for _, firstPair := range allPairAssetB {
		tmpFirst := strings.Replace(firstPair, asset, "", 1)
		var tmpArray []string
		// remove asset from third pair
		for _, thirdPair := range allPairAssetH {
			tmpSecond := strings.Replace(thirdPair, strings.ToLower(asset), "", 1)
			fExample := fmt.Sprintf("%s%s", strings.ToLower(tmpFirst), strings.ToLower(tmpSecond))
			sExample := fmt.Sprintf("%s%s", strings.ToLower(tmpSecond), strings.ToLower(tmpFirst))
			if commonfunction.FindElementArray(fExample, allPairH) {
				tmpArray = append(tmpArray, fmt.Sprintf("%s|%s", fExample, thirdPair))
			} else if commonfunction.FindElementArray(sExample, allPairH) {
				tmpArray = append(tmpArray, fmt.Sprintf("%s|%s", sExample, thirdPair))
			}
		}
		mapPairFinal[firstPair] = tmpArray
	}

	name_json := fmt.Sprintf("data/databinance/%s/%s_pairB_pairH_pairH.json", asset, asset)
	jsonStr, err := json.MarshalIndent(mapPairFinal, "", " ")
	if err != nil {
		log.Printf("Error: %s", err.Error())
	}

	_ = os.WriteFile(name_json, jsonStr, 0644)
}

func GetPairHuobi(asset string) {
	//get all pair for asset from Huobi

	respH, err := http.Get(urlHuobiExchangeInfo)
	if err != nil {
		panic(err)
	}
	defer respH.Body.Close()

	bodyH, _ := io.ReadAll(respH.Body)

	symbolsH := getinfohuobi.PairHuobi{}

	if err := json.Unmarshal(bodyH, &symbolsH); err != nil {
		panic(err)
	}

	allPairAssetH := []string{}
	allPairH := []string{}

	for _, valueH := range symbolsH.Data {
		allPairH = append(allPairH, valueH.Symbol)
		if valueH.BaseCurrency == strings.ToLower(asset) || valueH.QuoteCurrency == strings.ToLower(asset) {
			allPairAssetH = append(allPairAssetH, valueH.Symbol)
		}
	}

	// forming map for pair-pair-pair
	mapPairFinal := make(map[string][]string)

	// remove asset from first pair
	for _, firstPair := range allPairAssetH {
		tmpFirst := strings.Replace(firstPair, strings.ToLower(asset), "", 1)
		var tmpArray []string
		// remove asset from third pair
		for _, thirdPair := range allPairAssetH {
			tmpSecond := strings.Replace(thirdPair, strings.ToLower(asset), "", 1)
			fExample := fmt.Sprintf("%s%s", strings.ToLower(tmpFirst), strings.ToLower(tmpSecond))
			sExample := fmt.Sprintf("%s%s", strings.ToLower(tmpSecond), strings.ToLower(tmpFirst))
			if commonfunction.FindElementArray(fExample, allPairH) {
				tmpArray = append(tmpArray, fmt.Sprintf("%s|%s", fExample, thirdPair))
			} else if commonfunction.FindElementArray(sExample, allPairH) {
				tmpArray = append(tmpArray, fmt.Sprintf("%s|%s", sExample, thirdPair))
			}
		}
		mapPairFinal[firstPair] = tmpArray
	}

	name_json := fmt.Sprintf("data/datahuobi/%s/%s_pairH_pairH_pairH.json", asset, asset)
	jsonStr, err := json.MarshalIndent(mapPairFinal, "", " ")
	if err != nil {
		log.Printf("Error: %s", err.Error())
	}

	_ = os.WriteFile(name_json, jsonStr, 0644)
}

// GetPairOKX get and save from OKX pair-pair-pair
func GetPairOKX(asset string) {
	url := fmt.Sprintf("https://www.okx.com/priapi/v5/public/simpleProduct?instType=SPOT")

	resp, err := http.Get(url)
	if err != nil {
		log.Println("trable with get response from OKX", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Can't read body", err)
	}

	var pair getinfookx.PairOKX

	json.Unmarshal(body, &pair)

	allPairO := make([]string, len(pair.Data))
	var allPairAssetO []string

	// get all pair and pair for asset
	for i, j := range pair.Data {
		allPairO[i] = j.InstID
		tmp := strings.Split(j.InstID, "-")
		if tmp[0] == asset || tmp[1] == asset {
			allPairAssetO = append(allPairAssetO, allPairO[i])
		}
	}

	// forming map for pair-pair-pair
	mapPairFinalOKX := make(map[string][]string)

	// remove asset from first pair
	for _, firstPair := range allPairAssetO {
		tmpFirst := strings.Replace(firstPair, asset, "", 1)
		tmpFirst = strings.Replace(tmpFirst, "-", "", 1)
		var tmpArray []string
		// remove asset from third pair
		for _, thirdPair := range allPairAssetO {
			tmpSecond := strings.Replace(thirdPair, asset, "", 1)
			tmpSecond = strings.Replace(tmpSecond, "-", "", 1)
			fExample := fmt.Sprintf("%s-%s", tmpFirst, tmpSecond)
			sExample := fmt.Sprintf("%s-%s", tmpSecond, tmpFirst)
			if commonfunction.FindElementArray(fExample, allPairO) {
				tmpArray = append(tmpArray, fmt.Sprintf("%s|%s", fExample, thirdPair))
			} else if commonfunction.FindElementArray(sExample, allPairO) {
				tmpArray = append(tmpArray, fmt.Sprintf("%s|%s", sExample, thirdPair))
			}
		}
		mapPairFinalOKX[firstPair] = tmpArray
	}

	name_json := fmt.Sprintf("data/dataokx/%s/%s_pairO_pairO_pair).json", asset, asset)
	jsonStr, err := json.MarshalIndent(mapPairFinalOKX, "", " ")
	if err != nil {
		log.Printf("Error: %s", err.Error())
	}

	_ = os.WriteFile(name_json, jsonStr, 0644)
}
