package getinfohuobi

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

// GetListSymbolsAssetHuobi pair only for crypto
func GetListSymbolsAssetHuobi(asset string) ([]string, error) {

	//get all pair
	url := "https://api.huobi.pro/v1/common/symbols"

	resp, err := http.Get(url)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	symbols := PairHuobi{}

	if err := json.Unmarshal(body, &symbols); err != nil {
		panic(err)
	}

	allpair := []string{}

	for _, value := range symbols.Data {
		if strings.ToUpper(value.QuoteCurrency) == asset || strings.ToUpper(value.BaseCurrency) == asset {
			allpair = append(allpair, strings.ToUpper(value.Symbol))
		}
	}

	return allpair, nil

}

func GetInfoCryptoHuobi(fiat string) {
	url := "https://otc-cf.huobi.com/v1/otc/coin/config-trade?side=1&blockType=1"
	resp, err := http.Get(url)
	if err != nil {
		log.Println("trable with get response", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Can't read body", err)
	}

	var cf CryptoFiat

	json.Unmarshal(body, &cf)

	crypto, coinId := GetFiat(cf)

	coinId_json := fmt.Sprintf("data/datahuobi/%s/%s_coinId.json", fiat, fiat)
	jsonStr, err := json.MarshalIndent(coinId, "", " ")
	if err != nil {
		fmt.Printf("Error: %s", err.Error())
	}
	_ = os.WriteFile(coinId_json, jsonStr, 0644)

	getPair(fiat, crypto)
}

// key fiat, value - crypto
func GetFiat(cf CryptoFiat) (map[string][]string, map[string]int64) {
	fiatcryptaarray := make(map[string][]string)
	coinId := make(map[string]int64)

	for i, _ := range cf.Data {
		coinId[cf.Data[i].CryptoAsset.Name] = cf.Data[i].CryptoAsset.CoinID
		for _, l := range cf.Data[i].QuoteAsset {
			coinId[l.Name] = l.CoinID
			fiatcryptaarray[l.Name] = append(fiatcryptaarray[l.Name], cf.Data[i].CryptoAsset.Name)
		}
	}
	return fiatcryptaarray, coinId
}

func getPair(fiat string, crypto map[string][]string) {

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

	var ph PairHuobi

	json.Unmarshal(body, &ph)

	//list all pair on huobi
	var allpair []string
	for _, i := range ph.Data {
		allpair = append(allpair, fmt.Sprintf("%s%s", i.BaseCurrency, i.QuoteCurrency))
	}
	finalpair := make(map[string][]string)
	assets := crypto[fiat]

	for _, f := range assets {
		tmp := strings.ToLower(f)
		for _, s := range assets {
			s = strings.ToLower(s)
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

	name_json := fmt.Sprintf("data/datahuobi/%s/%s_pair.json", fiat, fiat)
	jsonStr, err := json.MarshalIndent(finalpair, "", " ")
	if err != nil {
		log.Printf("Error: %s", err.Error())
	}

	_ = os.WriteFile(name_json, jsonStr, 0644)
}
