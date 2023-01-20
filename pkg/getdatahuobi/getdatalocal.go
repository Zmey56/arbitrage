package getdatahuobi

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
)

func GetPairFromJSONHuobi(fiat string) map[string][]string {
	pair := ""
	switch fiat {
	case "AED":
		pair = fmt.Sprintf("data/datahuobi/%s/%s_pair.json", fiat, fiat)
	case "AMD":
		pair = fmt.Sprintf("data/datahuobi/%s/%s_pair.json", fiat, fiat)
	case "ARS":
		pair = fmt.Sprintf("data/datahuobi/%s/%s_pair.json", fiat, fiat)
	case "AZN":
		pair = fmt.Sprintf("data/datahuobi/%s/%s_pair.json", fiat, fiat)
	case "EUR":
		pair = fmt.Sprintf("data/datahuobi/%s/%s_pair.json", fiat, fiat)
	case "GEL":
		pair = fmt.Sprintf("data/datahuobi/%s/%s_pair.json", fiat, fiat)
	case "KZT":
		pair = fmt.Sprintf("data/datahuobi/%s/%s_pair.json", fiat, fiat)
	case "RUB":
		pair = fmt.Sprintf("data/datahuobi/%s/%s_pair.json", fiat, fiat)
	case "TRY":
		pair = fmt.Sprintf("data/datahuobi/%s/%s_pair.json", fiat, fiat)
	case "UAH":
		pair = fmt.Sprintf("data/datahuobi/%s/%s_pair.json", fiat, fiat)
	case "UZS":
		pair = fmt.Sprintf("data/datahuobi/%s/%s_pair.json", fiat, fiat)
	case "VND":
		pair = fmt.Sprintf("data/datahuobi/%s/%s_pair.json", fiat, fiat)
	default:
		log.Printf("For %v don't have para\n", fiat)
	}

	jsonfile, err := os.ReadFile(pair)
	if err != nil {
		panic(err)
	}
	var result map[string][]string
	_ = json.Unmarshal(jsonfile, &result)

	return result
}

func GetCurrencyHuobi(fiat string) []string {
	pathcurrency := fmt.Sprintf("data/datahuobi/%s/%s_pair.json", fiat, fiat)
	file, _ := os.Open(pathcurrency)
	defer file.Close()
	decoder := json.NewDecoder(file)
	var data map[string][]string
	decoder.Decode(&data)
	currency := []string{}
	for i, _ := range data {
		currency = append(currency, strings.ToUpper(i))
	}
	return currency
}
