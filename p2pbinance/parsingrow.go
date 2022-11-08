// Function for parsing data from Binance and get number pages

package p2pbinance

import (
	"encoding/json"
	"fmt"
	"io"
)

type Trader struct {
	nickName               string
	monthOrderCount        float64
	monthFinishRate        float64
	price                  string
	tradableQuantity       string
	maxSingleTransAmount   string
	minSingleTransAmount   string
	minSingleTransQuantity string
	tradeMethodName        string
	tradeType              string
	userNo                 string
	fiatUnit               string
	asset                  string
	tradeMethods           TradeMethods
}

type TradeMethods struct {
	identifier           string
	tradeMethodName      string
	tradeMethodShortName string
}

// All orders
type AllOrders []Trader

func ParsingJson(r io.Reader) (AllOrders, int) {
	allorders := AllOrders{}

	var result map[string]any

	body, _ := io.ReadAll(r)
	json.Unmarshal([]byte(body), &result)

	//numbers of rows
	numberRows := int(result["total"].(float64))

	for key, value := range result["data"].([]interface{}) {
		fmt.Println(key)
		for _, v := range value.(map[string]any) {
			for m, k := range v.(map[string]any) {
				order := Trader{}
				switch m {
				case "nickName":
					order.nickName = k.(string)
				case "monthOrderCount":
					order.monthOrderCount = k.(float64)
				case "monthFinishRate":
					order.monthFinishRate = k.(float64)
				case "price":
					order.price = k.(string)
				case "tradableQuantity":
					order.tradableQuantity = k.(string)
				case "maxSingleTransAmount":
					order.maxSingleTransAmount = k.(string)
				case "maxSingleTransQuantity":
					order.minSingleTransQuantity = k.(string)
				case "minSingleTransQuantity":
					order.minSingleTransQuantity = k.(string)
				case "tradeMethodName":
					order.tradeMethodName = k.(string)
				case "userNo":
					order.userNo = k.(string)
				case "fiatUnit":
					order.fiatUnit = k.(string)
				case "asset":
					order.asset = k.(string)
				case "tradeMethods":
					for _, d := range k.([]interface{}) {
						for i, j := range d.(map[string]any) {
							switch i {
							case "identifier":
								order.tradeMethods.identifier = j.(string)
							case "tradeMethodName":
								order.tradeMethods.tradeMethodName = j.(string)
							case "tradeMethodShortName":
								order.tradeMethods.tradeMethodShortName = j.(string)
							default:
								continue
							}
						}
					}
				default:
					continue

				}
				allorders = append(allorders, order)
			}
		}
	}
	numberPages := numberRows / 10

	return allorders, numberPages
}

//if namberPages > 1 {
//	//fmt.Println(string(jsonData))
//	for i := 2; i <= namberPages; i++ {
//		var m map[string]interface{}
//		err := json.Unmarshal(jsonData, &m)
//		if err != nil {
//			log.Println("Error", err)
//		}
//		m["page"] = i
//		newData, err := json.Marshal(m)
//		newJsonData := []byte(string(newData))
//		request, error := http.NewRequest("POST", httpposturl, bytes.NewBuffer(newJsonData))
//		request.Header.Set("Content-Type", "application/json; charset=UTF-8")
//
//		client := &http.Client{}
//		response, error := client.Do(request)
//		if error != nil {
//			panic(error)
//		}
//		defer response.Body.Close()
//
//		fmt.Println("response Status:", response.Status)
//		fmt.Println("response Headers:", response.Header)
//		body, _ := io.ReadAll(response.Body)
//		fmt.Println("response Body:", string(body))
//		if error != nil {
//			log.Println("Error", err)
//		}
//	}
//}
