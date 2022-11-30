//Get information from Binance about payment methods in dependes fiat

package getinfobinance

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

type Payment struct {
	Identifier           string
	PayAccount           string
	PayMethodId          string
	PayType              string
	TradeMethodName      string
	TradeMethodShortName string
}

var payments []Payment

func GetPeymontMethods(fiat ...string) []Payment {
	allpayment := payments
	httpposturl := "https://p2p.binance.com/bapi/c2c/v2/public/c2c/adv/filter-conditions"
	//fmt.Println("HTTP JSON POST URL:", httpposturl)

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
		panic(error)
	}
	defer response.Body.Close()

	body, _ := io.ReadAll(response.Body)

	var result map[string]any

	json.Unmarshal([]byte(body), &result)

	for key, value := range result {
		if key == "data" {
			for i, j := range value.(map[string]interface{}) {
				if i == "tradeMethods" {
					for _, m := range j.([]interface{}) {
						p := Payment{}
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
									p.PayMethodId = y.(string)
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
	return allpayment
}

func SavePaymentToJSON(p []Payment) {
	file, err := json.MarshalIndent(p, "", " ")
	if err != nil {
		panic(err)
	}
	os.WriteFile("data/rub_payment.json", file, 0644)
}

func UpdatePaymentToJSON(fiat ...string) {
	payment := GetPeymontMethods(fiat...)
	SavePaymentToJSON(payment)
	fmt.Println("Payment methods updated")
}

func UnpackedJSONOayments(p string) {

}
