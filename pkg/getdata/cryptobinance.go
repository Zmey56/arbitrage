// get all symbols from binance
package getdata

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/Zmey56/arbitrage/pkg/commonfunction"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

type exchangeInfo struct {
	symbol     string
	baseAsset  string
	quoteAsset string
}

// GetListSymbolsAssetBinance pair only for crypto
func GetListSymbolsAssetBinance(asset string) ([]string, error) {

	//get all pair
	url := "https://api.binance.com/api/v3/exchangeInfo"

	resp, err := http.Get(url)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	symbols := Symbols{}

	if err := json.Unmarshal(body, &symbols); err != nil {
		panic(err)
	}

	allpair := []string{}

	for _, value := range symbols.Symbols {
		if value.BaseAsset == asset || value.QuoteAsset == asset {
			allpair = append(allpair, strings.ToUpper(value.Symbol))
		}
	}

	return allpair, nil

}

func GetListSymbolsBinance(fiat string) {
	mapassets := GetAssets(fiat)
	assets := make([]string, len(mapassets))

	i := 0
	for asset, _ := range mapassets {
		assets[i] = asset
		i++
	}

	url := "https://api.binance.com/api/v3/exchangeInfo"

	resp, err := http.Get(url)

	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var dat map[string]interface{}

	if err := json.Unmarshal(body, &dat); err != nil {
		panic(err)
	}

	var arrayinfo []exchangeInfo

	var allpair []string

	for key, value := range dat {
		if key == "symbols" {
			for _, j := range value.([]interface{}) {
				exchangeInfo := exchangeInfo{}
				for i, m := range j.(map[string]interface{}) {
					switch i {
					case "symbol":
						exchangeInfo.symbol = m.(string)
						allpair = append(allpair, m.(string))
					case "baseAsset":
						exchangeInfo.baseAsset = m.(string)
					case "quoteAsset":
						exchangeInfo.quoteAsset = m.(string)
					default:
						continue
					}
				}
				arrayinfo = append(arrayinfo, exchangeInfo)
			}
		}
	}

	finalpair := make(map[string][]string)

	for _, f := range assets {
		tmp := f
		for _, s := range assets {
			if s != tmp {
				tmp_b := tmp + s
				tmp_q := s + tmp
				for _, pair := range allpair {
					if pair == tmp_b || pair == tmp_q {
						finalpair[tmp] = append(finalpair[tmp], pair)
					}
				}
			} else {
				continue
			}
		}
	}

	name_json := fmt.Sprintf("data/databinance/%s/%s_pair.json", fiat, fiat)
	jsonStr, err := json.MarshalIndent(finalpair, "", " ")
	if err != nil {
		log.Printf("Error: %s", err.Error())
	}

	_ = os.WriteFile(name_json, jsonStr, 0644)

}

func GetAssets(fiat ...string) map[string]string {

	assets := make(map[string]string)
	httpposturl := "https://p2p.binance.com/bapi/c2c/v2/friendly/c2c/portal/config"

	var jsonData = []byte(`{
		"fiat": "USD"
	}`)
	if fiat != nil {
		var m map[string]interface{}
		err := json.Unmarshal(jsonData, &m)
		if err != nil {
			log.Println("Error", err)
		}

		m["fiat"] = fiat[0]
		newData, err := json.Marshal(m)
		jsonData = []byte(string(newData))
	}

	request, error := http.NewRequest("POST", httpposturl, bytes.NewBuffer(jsonData))
	request.Header.Set("Content-Type", "application/json; charset=UTF-8")

	client := &http.Client{}
	response, error := client.Do(request)
	if error != nil {
		log.Panic(error)
		os.Exit(1)
	}
	defer response.Body.Close()

	body, _ := io.ReadAll(response.Body)

	var result map[string]any

	json.Unmarshal([]byte(body), &result)

	for key, value := range result {
		if key == "data" {
			for i, j := range value.(map[string]interface{}) {
				if i == "areas" {
					for _, m := range j.([]interface{}) {
						for l, k := range m.(map[string]interface{}) {
							if l == "tradeSides" {
								for _, p := range k.([]interface{}) {
									for t, u := range p.(map[string]interface{}) {
										if t == "assets" {
											for _, a := range u.([]interface{}) {
												key := a.(map[string]interface{})["asset"]
												value := a.(map[string]interface{})["description"]
												if key != nil && value != nil {
													k := key.(string)
													v := value.(string)
													assets[k] = v
												}
											}
										}
									}
								}
							}
						}
					}
				}
			}
		}
	}
	log.Println("assets", assets)
	return assets
}

func GetAssetsLocalBinance(fiat string) []string {
	pathcurrency := fmt.Sprintf("data/databinance/%s/%s_pair.json", fiat, fiat)
	file, _ := os.Open(pathcurrency)
	defer file.Close()
	decoder := json.NewDecoder(file)
	var data map[string][]string
	decoder.Decode(&data)
	assets := []string{}
	for i, _ := range data {
		assets = append(assets, strings.ToUpper(i))
	}
	return assets
}

// GetListSymbolsBinancePairPair - get for working currency - p2p - pair - pair - currency
func GetListSymbolsBinancePairPair(fiat string) {
	//get slice all assets for this fiat
	mapassets := GetAssets(fiat)
	assets := make([]string, len(mapassets))

	i := 0
	for asset, _ := range mapassets {
		assets[i] = asset
		i++
	}

	//get all pair
	url := "https://api.binance.com/api/v3/exchangeInfo"

	resp, err := http.Get(url)

	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	symbols := Symbols{}

	if err := json.Unmarshal(body, &symbols); err != nil {
		panic(err)
	}

	allpair := make([]string, len(symbols.Symbols))

	finalpair := make(map[string][]string)

	rubPair := []string{}

	for i, value := range symbols.Symbols {
		allpair[i] = value.Symbol
	}

	for _, value := range symbols.Symbols {
		if value.QuoteAsset == fiat || value.BaseAsset == fiat {
			rubPair = append(rubPair, value.Symbol)

			tmpNotFiat := ""
			if value.QuoteAsset != fiat {
				tmpNotFiat = value.QuoteAsset
			} else {
				tmpNotFiat = value.BaseAsset
			}
			log.Println("tmpNotFiat", tmpNotFiat)
			for _, a := range assets {
				tmp_b := tmpNotFiat + a
				tmp_q := a + tmpNotFiat
				//log.Println(tmp_b, tmp_q)
				for _, pair := range allpair {
					if tmp_b == pair {
						finalpair[a] = append(finalpair[a], fmt.Sprintf("%s|%s", tmp_b, value.Symbol))
					}
					if tmp_q == pair {
						finalpair[a] = append(finalpair[a], fmt.Sprintf("%s|%s", tmp_q, value.Symbol))
					}
				}
			}
		}
	}

	name_json := fmt.Sprintf("data/databinance/%s/%s_pair_pair.json", fiat, fiat)
	jsonStr, err := json.MarshalIndent(finalpair, "", " ")
	if err != nil {
		log.Printf("Error: %s", err.Error())
	}

	_ = os.WriteFile(name_json, jsonStr, 0644)

}

// GetThreePairs - get for working pair - pair - pair
func GetThreePairs(asset string) {
	//get slice all pair for asset for this fiat
	pairs, _ := GetListSymbolsAssetBinance(asset)

	mapPairs := make(map[string][]string)

	for i, pair := range pairs {
		tmpPairs := make([]string, 0, len(pairs)-1)
		for j, p := range pairs {
			if i != j {
				tmpPairs = append(tmpPairs, p)
			}
		}
		mapPairs[pair] = tmpPairs
	}

	mapPairFinal := make(map[string][]string)
	allPair := GetAllPair()

	for firstPair, arrayPair := range mapPairs {
		tmpFirst := strings.Replace(firstPair, "USDT", "", 1)
		var tmpArray []string
		for _, secondPair := range arrayPair {
			tmpSecond := strings.Replace(secondPair, "USDT", "", 1)
			fExample := fmt.Sprintf("%s%s", tmpFirst, tmpSecond)
			sExample := fmt.Sprintf("%s%s", tmpSecond, tmpFirst)
			if commonfunction.FindElementArray(fExample, allPair) {
				tmpArray = append(tmpArray, fmt.Sprintf("%s|%s", fExample, secondPair))
			} else if commonfunction.FindElementArray(sExample, allPair) {
				tmpArray = append(tmpArray, fmt.Sprintf("%s|%s", sExample, secondPair))
			}
		}
		mapPairFinal[firstPair] = uniqElement(tmpArray)
	}

	name_json := fmt.Sprintf("data/databinance/%s/%s_pair_pair_pair.json", asset, asset)
	jsonStr, err := json.MarshalIndent(mapPairFinal, "", " ")
	if err != nil {
		log.Printf("Error: %s", err.Error())
	}

	_ = os.WriteFile(name_json, jsonStr, 0644)
}

func findAndRemove(pair string, pairs []string) []string {
	var index int
	for i, p := range pairs {
		if p == pair {
			index = i
			break
		}
	}

	// Remove the element at the index
	pairs = append(pairs[:index], pairs[index+1:]...)

	return pairs
}

// GetAllPair get all pair in Binance
func GetAllPair() []string {
	url := "https://api.binance.com/api/v3/exchangeInfo"

	resp, err := http.Get(url)

	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	symbols := Symbols{}

	if err := json.Unmarshal(body, &symbols); err != nil {
		panic(err)
	}

	allpair := make([]string, len(symbols.Symbols))

	for i, value := range symbols.Symbols {
		allpair[i] = value.Symbol
	}

	return allpair
}

func uniqElement(a []string) []string {
	uniqueE := make(map[string]bool)
	for _, elem := range a {
		uniqueE[elem] = true
	}
	a = []string{}
	for price := range uniqueE {
		a = append(a, price)
	}
	return a
}

func GetPairWithoutB(fiat string) {
	//get all pair
	url := "https://api.binance.com/api/v3/exchangeInfo"

	resp, err := http.Get(url)
	if err != nil {
		log.Println(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
	}

	symbols := Symbols{}

	if err := json.Unmarshal(body, &symbols); err != nil {
		panic(err)
	}

	var fiatpair []string
	fiatLow := strings.ToLower(fiat)
	for _, i := range symbols.Symbols {
		if strings.ToLower(i.BaseAsset) == fiatLow || strings.ToLower(i.QuoteAsset) == fiatLow {
			fiatpair = append(fiatpair, strings.Join([]string{i.BaseAsset, i.QuoteAsset}, ""))
		}
	}

	name_json := fmt.Sprintf("data/databinance/%s/%s_pair_without.json", fiat, fiat)
	jsonStr, err := json.MarshalIndent(fiatpair, "", " ")
	if err != nil {
		log.Printf("Error: %s", err.Error())
	}

	_ = os.WriteFile(name_json, jsonStr, 0644)

}
