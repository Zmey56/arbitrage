// get all symbols from binance
package getinfobinance

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

type exchangeInfo struct {
	symbol     string
	baseAsset  string
	quoteAsset string
}

func GetListSymbols(mapassets map[string]string, fiat string) {
	assets := make([]string, len(mapassets))

	i := 0
	for asset, _ := range mapassets {
		assets[i] = asset
		i++
	}
	//fmt.Println(assets, len(assets))
	//

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
	t := time.Now()
	t.String()
	name_json := fmt.Sprintf("data/%s_pair_%s.json", fiat, t.Format("2006_01_02"))
	jsonStr, err := json.MarshalIndent(finalpair, "", " ")
	if err != nil {
		fmt.Printf("Error: %s", err.Error())
	}

	_ = os.WriteFile(name_json, jsonStr, 0644)

}
