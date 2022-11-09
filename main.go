package main

import (
	"github.com/Zmey56/arbitrage/getinfobinance"
)

func main() {
	getinfobinance.GetAssets("GEL")

	//payTypes := []string{}
	//
	//jsonData := p2pbinance.GetDataP2P("USDT", "RUB", "Buy", payTypes)
	//
	//httpposturl := "https://p2p.binance.com/bapi/c2c/v2/friendly/c2c/adv/search"
	//fmt.Println("HTTP JSON POST URL:", httpposturl)
	//
	//request, error := http.NewRequest("POST", httpposturl, bytes.NewBuffer(jsonData))
	//request.Header.Set("Content-Type", "application/json; charset=UTF-8")
	//
	//client := &http.Client{}
	//response, error := client.Do(request)
	//if error != nil {
	////	panic(error)
	////}
	////defer response.Body.Close()
	////
	//////tmp, _ := p2pbinance.ParsingJson(response.Body)
	//
	//jsonDatabank := p2pbinance.GetDataP2PBank("RUB")
	//
	//httpposturlbanks := "https://p2p.binance.com/bapi/c2c/v2/public/c2c/adv/filter-conditions"
	//
	//fmt.Println("HTTP JSON POST URL:", httpposturlbanks)
	//
	//request_bank, error := http.NewRequest("POST", httpposturlbanks, bytes.NewBuffer(jsonDatabank))
	//request_bank.Header.Set("Content-Type", "application/json; charset=UTF-8")
	//
	//client_bank := &http.Client{}
	//response_bank, error := client_bank.Do(request_bank)
	//if error != nil {
	//	panic(error)
	//}
	//fmt.Println(response_bank)
	//defer response_bank.Body.Close()
	//
	//jsonDataconfig := p2pbinance.GetDataP2PCrypto("RUB")
	//
	//httpposturlconfig := "https://p2p.binance.com/bapi/c2c/v2/friendly/c2c/portal/config"
	//
	//fmt.Println("HTTP JSON POST URL:", httpposturlbanks)
	//
	//request_config, error := http.NewRequest("POST", httpposturlconfig, bytes.NewBuffer(jsonDataconfig))
	//request_config.Header.Set("Content-Type", "application/json; charset=UTF-8")
	//
	//client_config := &http.Client{}
	//response_config, error := client_config.Do(request_config)
	//if error != nil {
	//	panic(error)
	//}
	//fmt.Println(response_config)
	//defer response_config.Body.Close()

}
