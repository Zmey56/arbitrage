//Get information from Binance about payment methods in dependes fiat

package getinfobinance

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
)

func GetPeymontMethodsBinance(fiat string) {
	allpayment := []PaymentBinance{}
	httpposturl := "https://p2p.binance.com/bapi/c2c/v2/public/c2c/adv/filter-conditions"

	m := make(map[string]interface{})

	m["fiat"] = fiat
	newData, _ := json.Marshal(m)
	jsonData := []byte(string(newData))

	request, error := http.NewRequest("POST", httpposturl, bytes.NewBuffer(jsonData))
	request.Header.Set("Content-Type", "application/json; charset=UTF-8")

	client := &http.Client{}
	response, error := client.Do(request)
	if error != nil {
		panic(error)
	}
	defer response.Body.Close()

	body, _ := io.ReadAll(response.Body)

	var result map[string]any

	json.Unmarshal([]byte(body), &result)

	//log.Println(result)

	for key, value := range result {
		if key == "data" {
			for i, j := range value.(map[string]interface{}) {
				if i == "tradeMethods" {
					for _, m := range j.([]interface{}) {
						p := PaymentBinance{}
						for x, y := range m.(map[string]interface{}) {

							switch x {
							case "identifier":
								if y != nil {
									p.Identifier = y.(string)
								}
							case "payAccount":
								if y != nil {
									p.PayAccount = y.(string)
								}
							case "payMethodId":
								if y != nil {
									switch y.(type) {
									case int:
										p.PayMethodId = y.(int)
									case string:
										p.PayMethodId, _ = strconv.Atoi(y.(string))
									}
								}
							case "payType":
								if y != nil {
									p.PayType = y.(string)
								}
							case "tradeMethodName":
								if y != nil {
									p.TradeMethodName = y.(string)
								}
							case "tradeMethodShortName":
								if y != nil {
									p.TradeMethodShortName = y.(string)
								}
							default:
								continue
							}
						}
						allpayment = append(allpayment, p)
					}
				}
			}
		}
	}
	savePaymentToJSONBinance(allpayment, fiat)
}

func savePaymentToJSONBinance(p []PaymentBinance, fiat string) {
	file, err := json.MarshalIndent(p, "", " ")
	if err != nil {
		panic(err)
	}
	path_json := fmt.Sprintf("data/databinance/%s/%s_payment.json", fiat, fiat)
	os.WriteFile(path_json, file, 0644)
}
