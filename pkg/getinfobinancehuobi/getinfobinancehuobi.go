package getinfobinancehuobi

import (
	"encoding/json"
	"fmt"
	"github.com/Zmey56/arbitrage/pkg/getdata"
	"github.com/Zmey56/arbitrage/pkg/getdatahuobi"
	"github.com/Zmey56/arbitrage/pkg/getinfohuobi"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

func GetPairBinanceHuobiHuobi(fiat string) {
	//Map for assets from binance and coin from Huobi
	assetscoinpair := make(map[string][]string)

	//Get Assets from Binance
	assets := getdata.GetAssets(fiat)
	assets_symbol := make([]string, 0, len(assets))

	for k, _ := range assets {
		assets_symbol = append(assets_symbol, k)
	}

	//Get Coin from Huobi
	currencyarr := getdatahuobi.GetCurrencyHuobi(fiat)

	//find pair from Huobi

	url := "https://api.huobi.pro/v1/common/symbols"
	resp, err := http.Get(url)
	if err != nil {
		log.Println("trable with get response for pair", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Can't read body for pair", err)
	}

	var ph getinfohuobi.PairHuobi

	json.Unmarshal(body, &ph)

	for _, f := range assets_symbol {
		pair := []string{}
		for _, i := range ph.Data {
			for _, j := range currencyarr {
				if i.BaseCurrency == strings.ToLower(f) && i.QuoteCurrency == strings.ToLower(j) {
					pair = append(pair, i.Symbol)
				}
				if i.QuoteCurrency == strings.ToLower(f) && i.BaseCurrency == strings.ToLower(j) {
					pair = append(pair, i.Symbol)
				}
			}

		}
		assetscoinpair[f] = pair
	}

	name_json := fmt.Sprintf("data/databinancehuobi/%s/%s_pair.json", fiat, fiat)
	jsonStr, err := json.MarshalIndent(assetscoinpair, "", " ")
	if err != nil {
		log.Printf("Error: %s", err.Error())
	}

	_ = os.WriteFile(name_json, jsonStr, 0644)
}

func GetPairHuobiHuobiBinance(fiat string) {
	//Map for assets from binance and coin from Huobi
	assetscoinpair := make(map[string][]string)

	//Get Coin from Huobi
	coinarr := getdatahuobi.GetCurrencyHuobi(fiat)

	//Get Assets from Huobi
	assets := getdata.GetAssets(fiat)
	assets_symbol := make([]string, 0, len(assets))

	for k, _ := range assets {
		assets_symbol = append(assets_symbol, k)
	}

	//find pair from Huobi

	url := "https://api.huobi.pro/v1/common/symbols"
	resp, err := http.Get(url)
	if err != nil {
		log.Println("trable with get response for pair", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Can't read body for pair", err)
	}

	var ph getinfohuobi.PairHuobi

	json.Unmarshal(body, &ph)

	for _, f := range coinarr {
		pair := []string{}
		for _, i := range ph.Data {
			for _, j := range assets_symbol {
				if i.BaseCurrency == strings.ToLower(f) && i.QuoteCurrency == strings.ToLower(j) {
					pair = append(pair, i.Symbol)
				}
				if i.QuoteCurrency == strings.ToLower(f) && i.BaseCurrency == strings.ToLower(j) {
					pair = append(pair, i.Symbol)
				}
			}

		}
		assetscoinpair[f] = pair
	}
	//log.Println("RESULT", assetscoinpair)

	name_json := fmt.Sprintf("data/datahuobibinance/%s/%s_pair.json", fiat, fiat)
	jsonStr, err := json.MarshalIndent(assetscoinpair, "", " ")
	if err != nil {
		log.Printf("Error: %s", err.Error())
	}

	_ = os.WriteFile(name_json, jsonStr, 0644)
}
