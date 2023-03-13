// get all symbols from binance
package getdata

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

type exchangeInfo struct {
	symbol     string
	baseAsset  string
	quoteAsset string
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
