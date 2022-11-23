// get all symbols from binance
package getinfobinance

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type exchangeInfo struct {
	symbol     string
	baseAsset  string
	quoteAsset string
}

func GetListSymbols(asset string) []exchangeInfo {
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

	for key, value := range dat {
		if key == "symbols" {
			for _, j := range value.([]interface{}) {
				exchangeInfo := exchangeInfo{}
				for i, m := range j.(map[string]interface{}) {
					switch i {
					case "symbol":
						exchangeInfo.symbol = m.(string)
						fmt.Println("struct", exchangeInfo.symbol)
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
	return arrayinfo
}
