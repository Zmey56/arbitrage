//in this function get available crypto for fiat fot working

package getinfobinance

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

func GetAssets(fiat ...string) map[string]string {
	assets := make(map[string]string)
	httpposturl := "https://p2p.binance.com/bapi/c2c/v2/friendly/c2c/portal/config"
	fmt.Println("HTTP JSON POST URL:", httpposturl)

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

	fmt.Println("response Status:", response.Status)
	fmt.Println("response Headers:", response.Header)
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
	return assets
}
