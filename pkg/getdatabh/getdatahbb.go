package getdatabh

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

func GetPairFromJSONHBB(fiat string) map[string][]string {
	pair := ""
	switch fiat {
	case "AED":
		pair = fmt.Sprintf("data/datahbb/%s/%s_pair.json", fiat, fiat)
	case "AMD":
		pair = fmt.Sprintf("data/datahbb/%s/%s_pair.json", fiat, fiat)
	case "ARS":
		pair = fmt.Sprintf("data/datahbb/%s/%s_pair.json", fiat, fiat)
	case "AZN":
		pair = fmt.Sprintf("data/datahbb/%s/%s_pair.json", fiat, fiat)
	case "EUR":
		pair = fmt.Sprintf("data/datahbb/%s/%s_pair.json", fiat, fiat)
	case "GEL":
		pair = fmt.Sprintf("data/datahbb/%s/%s_pair.json", fiat, fiat)
	case "KZT":
		pair = fmt.Sprintf("data/datahbb/%s/%s_pair.json", fiat, fiat)
	case "RUB":
		pair = fmt.Sprintf("data/datahbb/%s/%s_pair.json", fiat, fiat)
	case "TRY":
		pair = fmt.Sprintf("data/datahbb/%s/%s_pair.json", fiat, fiat)
	case "UAH":
		pair = fmt.Sprintf("data/datahbb/%s/%s_pair.json", fiat, fiat)
	case "USD":
		pair = fmt.Sprintf("data/datahbb/%s/%s_pair.json", fiat, fiat)
	case "UZS":
		pair = fmt.Sprintf("data/datahbb/%s/%s_pair.json", fiat, fiat)
	case "VND":
		pair = fmt.Sprintf("data/datahbb/%s/%s_pair.json", fiat, fiat)
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
